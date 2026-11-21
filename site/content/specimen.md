---
title: Backpropagation is the chain rule; the engineering is everything else
description: Reverse-mode automatic differentiation in forty lines of Go, and where the memory goes.
kind: essay
noindex: true
date: "2027-03"
abstract: Reverse-mode automatic differentiation fits in forty lines of Go. Making it fast, and making it fit in memory, takes everything a systems engineer already knows.
---

<!-- Typesetting specimen: exercises every component in DESIGN.md (MathML, footnotes, code, a table,
an algorithm, exercises). Not linked from the site; used to check the typography after changes. -->

## <span class="n">I.</span>The rule

Every derivative a neural network needs comes from one fact. If <math><mi>y</mi><mo>=</mo><mi>f</mi><mo form="prefix" stretchy="false" lspace="0" rspace="0">(</mo><mi>u</mi><mo form="postfix" stretchy="false" lspace="0" rspace="0">)</mo></math> and <math><mi>u</mi><mo>=</mo><mi>g</mi><mo form="prefix" stretchy="false" lspace="0" rspace="0">(</mo><mi>x</mi><mo form="postfix" stretchy="false" lspace="0" rspace="0">)</mo></math>, the derivative of <math><mi>y</mi></math> with respect to <math><mi>x</mi></math> is the product of the two local derivatives. Backpropagation applies that fact to a graph instead of a chain: each node <math><msub><mi>v</mi><mi>i</mi></msub></math> receives an *adjoint*, the derivative of the loss with respect to that node, and passes a share of it back to each of its inputs:[^griewank]

<div class="eq num" id="eq-1"><math display="block"><msub><mover accent="false"><mi>v</mi><mo>&#xAF;</mo></mover><mi>j</mi></msub><mo>=</mo><munder><mo>&sum;</mo><mrow><mi>i</mi><mo>&isin;</mo><mi>out</mi><mo form="prefix" stretchy="false" lspace="0" rspace="0">(</mo><mi>j</mi><mo form="postfix" stretchy="false" lspace="0" rspace="0">)</mo></mrow></munder><msub><mover accent="false"><mi>v</mi><mo>&#xAF;</mo></mover><mi>i</mi></msub><mspace width="0.2em"/><mfrac><mrow><mo>&part;</mo><msub><mi>v</mi><mi>i</mi></msub></mrow><mrow><mo>&part;</mo><msub><mi>v</mi><mi>j</mi></msub></mrow></mfrac><mo>.</mo></math><span class="no">(1)</span></div>

Equation&nbsp;(1) becomes an algorithm as soon as the nodes are visited in the right order.

<Algorithm id="algorithm-b" letter="B" title="Backpropagation">

Given nodes <math><msub><mi>v</mi><mn>1</mn></msub><mo>,</mo><mo>&hellip;</mo><mo>,</mo><msub><mi>v</mi><mi>n</mi></msub></math> in topological order, with the loss at <math><msub><mi>v</mi><mi>n</mi></msub></math>, this algorithm computes every adjoint <math><msub><mover accent="false"><mi>v</mi><mo>&#xAF;</mo></mover><mi>j</mi></msub></math>.

<Steps>
<Step n="B1" label="Initialize.">

Set <math><msub><mover accent="false"><mi>v</mi><mo>&#xAF;</mo></mover><mi>n</mi></msub><mo>&larr;</mo><mn>1</mn></math> and <math><msub><mover accent="false"><mi>v</mi><mo>&#xAF;</mo></mover><mi>j</mi></msub><mo>&larr;</mo><mn>0</mn></math> for <math><mn>1</mn><mo>&le;</mo><mi>j</mi><mo>&lt;</mo><mi>n</mi></math>. Set <math><mi>i</mi><mo>&larr;</mo><mi>n</mi></math>.

</Step>
<Step n="B2" label="Propagate.">

For each input <math><msub><mi>v</mi><mi>j</mi></msub></math> of <math><msub><mi>v</mi><mi>i</mi></msub></math>, set <math><msub><mover accent="false"><mi>v</mi><mo>&#xAF;</mo></mover><mi>j</mi></msub><mo>&larr;</mo><msub><mover accent="false"><mi>v</mi><mo>&#xAF;</mo></mover><mi>j</mi></msub><mo>+</mo><msub><mover accent="false"><mi>v</mi><mo>&#xAF;</mo></mover><mi>i</mi></msub><mspace width="0.1em"/><mo>&part;</mo><msub><mi>v</mi><mi>i</mi></msub><mo lspace="0" rspace="0">/</mo><mo>&part;</mo><msub><mi>v</mi><mi>j</mi></msub></math>.

</Step>
<Step n="B3" label="Decrease *i*." end>

Decrease <math><mi>i</mi></math> by&nbsp;1. If <span class="nowrap"><math><mi>i</mi><mo>&gt;</mo><mn>0</mn></math>,</span> return to&nbsp;[B2](#step-b2); otherwise the algorithm terminates.

</Step>
</Steps>

</Algorithm>

## <span class="n">II.</span>A value that remembers

In Go, [step&nbsp;B2](#step-b2) is a closure stored on each value when the forward pass creates it.

```go
type Value struct {
	Data, Grad float64
	parents    []*Value
	backward   func()
}

// Mul records both inputs so step B2 can reach them.
func Mul(a, b *Value) *Value {
	out := &Value{
		Data:    a.Data * b.Data,
		parents: []*Value{a, b},
	}
	out.backward = func() {
		a.Grad += b.Data * out.Grad
		b.Grad += a.Data * out.Grad
	}
	return out
}
```

The `+=` is the sum in&nbsp;[(1)](#eq-1). Forgetting it, so that a node used twice keeps only its last contribution, is the classic first bug.

## <span class="n">III.</span>What the tutorials skip

The code allocates a node for every multiplication and keeps all of them alive until [step&nbsp;B3](#step-b3) finishes. At the scale of a real model, that is where the memory goes.[^adam]

| Buffer | Lives from | Until | Grows with |
|---|---|---|---|
| Parameters | start | end of run | model size |
| Optimizer state | start | end of run | 2&nbsp;&times;&nbsp;model size |
| Gradients | step&nbsp;B1 | optimizer step | model size |
| Activations | forward pass | step&nbsp;B3 | batch&nbsp;&times;&nbsp;depth |

<Exercises>
<Exercise n="1" rating="10">

What does Algorithm&nbsp;B compute for a node that no other node uses?

</Exercise>
<Exercise n="2" rating="M20" recommended>

Prove that when Algorithm&nbsp;B terminates, <math><msub><mover accent="false"><mi>v</mi><mo>&#xAF;</mo></mover><mi>j</mi></msub><mo>=</mo><mo>&part;</mo><msub><mi>v</mi><mi>n</mi></msub><mo lspace="0" rspace="0">/</mo><mo>&part;</mo><msub><mi>v</mi><mi>j</mi></msub></math> for <span class="nowrap">all&nbsp;<math><mi>j</mi></math>.</span>

</Exercise>
<Exercise n="3" rating="25">

Change `Mul` so that it needs no closure, and measure what that saves.

</Exercise>
</Exercises>

<Answers>
<Answer n="1">

Zero. No step adds anything to its adjoint, which is right: the loss doesn&rsquo;t depend on that node.

</Answer>
<Answer n="2">

Work downward from <math><mi>n</mi></math>. Every node that uses <math><msub><mi>v</mi><mi>j</mi></msub></math> comes later in topological order, so all of their adjoints are final before <math><msub><mover accent="false"><mi>v</mi><mo>&#xAF;</mo></mover><mi>j</mi></msub></math> is complete, and [step&nbsp;B2](#step-b2) has added one term for each of them. That sum is the multivariable chain rule.

</Answer>
<Answer n="3">

Store an operation code and the two inputs in the value itself, and give the type one `backward` method that switches on the code. That removes two heap allocations per node, the closure and the `parents` slice.

</Answer>
</Answers>

[^griewank]: Griewank and Walther give the full treatment in *Evaluating Derivatives* (2008).
[^adam]: Optimizer state assumes Adam, which keeps two moments per parameter.
