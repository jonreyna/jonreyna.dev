---
title: Backpropagation is the chain rule; the engineering is everything else
description: Reverse-mode automatic differentiation in forty lines of Go, and where the memory goes.
kind: essay
noindex: true
date: "2027-03"
abstract: Reverse-mode automatic differentiation fits in forty lines of Go. Making it fast, and making it fit in memory, takes everything a systems engineer already knows.
---

<!-- Typesetting specimen: exercises every component in DESIGN.md (TeX math, footnotes, code, a table,
an algorithm, exercises). Not linked from the site; used to check the typography after changes. -->

## The rule

Every derivative a neural network needs comes from one fact. If $y = f(u)$ and $u = g(x)$, the derivative of $y$ with respect to $x$ is the product of the two local derivatives. Backpropagation applies that fact to a graph instead of a chain: each node $v_i$ receives an *adjoint*, the derivative of the loss with respect to that node, and passes a share of it back to each of its inputs:[^griewank]

<div class="eq num" id="eq-1">$$\bar{v}_j = \sum_{i \in \mathrm{out}(j)} \bar{v}_i \, \frac{\partial v_i}{\partial v_j}.$$<span class="no">(1)</span></div>

Equation&nbsp;(1) becomes an algorithm as soon as the nodes are visited in the right order.

<Algorithm id="algorithm-b" letter="B" title="Backpropagation">

Given nodes $v_1, \ldots, v_n$ in topological order, with the loss at $v_n$, this algorithm computes every adjoint $\bar{v}_j$.

<Steps>
<Step n="B1" label="Initialize.">

Set $\bar{v}_n \gets 1$ and $\bar{v}_j \gets 0$ for $1 \le j < n$. Set $i \gets n$.

</Step>
<Step n="B2" label="Propagate.">

For each input $v_j$ of $v_i$, set $\bar{v}_j \gets \bar{v}_j + \bar{v}_i \, \partial v_i / \partial v_j$.

</Step>
<Step n="B3" label="Decrease *i*." end>

Decrease $i$ by&nbsp;1. If $i > 0$, return to&nbsp;[B2](#step-b2); otherwise the algorithm terminates.

</Step>
</Steps>

</Algorithm>

## A value that remembers

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

## What the tutorials skip

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

Prove that when Algorithm&nbsp;B terminates, $\bar{v}_j = \partial v_n / \partial v_j$ for <span class="nowrap">all&nbsp;$j$.</span>

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

Work downward from $n$. Every node that uses $v_j$ comes later in topological order, so all of their adjoints are final before $\bar{v}_j$ is complete, and [step&nbsp;B2](#step-b2) has added one term for each of them. That sum is the multivariable chain rule.

</Answer>
<Answer n="3">

Store an operation code and the two inputs in the value itself, and give the type one `backward` method that switches on the code. That removes two heap allocations per node, the closure and the `parents` slice.

</Answer>
</Answers>

[^griewank]: Griewank and Walther give the full treatment in *Evaluating Derivatives* (2008).
[^adam]: Optimizer state assumes Adam, which keeps two moments per parameter.
