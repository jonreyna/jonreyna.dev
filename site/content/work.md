---
title: Work
description: Selected work by Jon Reyna, from card authorization in C to GitOps for 2,317 Kubernetes clusters and orchestration for AI inference.
kind: page
---

Thirteen years of distributed systems, most of them keeping large fleets of infrastructure in their intended state. These are the pieces I would point to first, newest first; the [résumé](/resume) has every role.

<h2 id="equinix">Equinix, 2023–2026</h2>

<p class="work"><span class="t">Fabric Super Agent skills</span>, <span class="m">2026.</span> Connection-search and inventory skills for Equinix’s AI agent, built over its MCP server. The rewrite shipped in August 2026 to the Equinix Portal, Slack and Microsoft Teams, and a recurring production error went away with it. I also hardened the LLM-as-judge harness that evaluates the agent, so that a crash or a timeout in the judge no longer counts as a wrong answer.</p>
<p class="work"><span class="t">Fabric Application Connect</span>, <span class="m">2025.</span> A private connectivity service for AI inference. With one other engineer I owned its platform orchestration through the December 2025 production release: Temporal workflows for its core resources, integration with Network Edge virtual devices and Fabric Cloud Router, and the gRPC and REST boundary between the customer API and the workers.</p>
<p class="work"><span class="t">Temporal middleware</span>, <span class="m">2024–2025.</span> Retries, logging, Prometheus metrics and rollback for Temporal workflows and activities, written as middleware that stacks the way go-kit’s does. Getting wrapped functions to register under their own names took a trip through the Temporal SDK in a debugger. Every engineer on the team built on it, and the AI-inference team adopted it later.</p>
<p class="work"><span class="t">Multicloud orchestration</span>, <span class="m">2024–2025.</span> Temporal workflows that provision connectivity from Equinix Fabric Cloud Routers to AWS, Google Cloud and Azure. It reached production and was demoed at Google Cloud Next ’24. I led deprovisioning: workflows that undo every deployment step, in all three clouds.</p>

<h2 id="the-home-depot">The Home Depot, 2020–2023</h2>

<p class="work"><span class="t">küber</span>, <span class="m">2020–2023.</span> GitOps deployment for the in-store Kubernetes platform, 2,317 clusters in stores and distribution centers. It renders each deployable unit through Go templates and kustomize and commits the result to every target store’s repository, where Flux syncs it. I was hired to replace a failing proof of concept and rewrote it from scratch. It still carries about 15,000 deployments a week for 24 teams, 99.8% of them successful, and the engineer I mentored on it now maintains it.</p>
<p class="work"><span class="t">An in-memory filesystem</span>, <span class="m">2020.</span> küber clones, renders and commits without touching disk, through one in-memory filesystem that implements both go-git’s and kustomize’s filesystem interfaces. That made every step unit-testable, months before Go 1.16 gave the standard library a filesystem interface of its own.</p>
<p class="work"><span class="t">Watchdog</span>, <span class="m">2021.</span> A dead man’s switch for the fleet’s monitoring. Each cluster’s Prometheus sends a heartbeat every 20 minutes; a store that falls silent raises an alert, which catches the network outages that monitoring inside the cluster can’t report. Clusters register themselves, and it watches 2,232 production targets.</p>

<h2 id="cray">Cray, 2019–2020</h2>

<p class="work"><span class="t">Shasta fabric controller</span>, <span class="m">2019–2020.</span> Go microservices behind a configuration API for the high-speed network of Cray’s Shasta exascale supercomputers, writing changes directly into switch memory. The design keeps the configuration you intend apart from what each switch is running, after RFC 8342, and borrows versioned keys and watches from etcd. I also moved the existing C++ controller to Go.</p>

<h2 id="juniper">Juniper Networks, 2015–2019</h2>

<p class="work"><span class="t">Network automation in Go</span>, <span class="m">2015–2019.</span> On site at WOW!, a U.S. cable operator in 19 markets, whose engineers did nearly everything by logging into routers by hand. I introduced Go and built the automation platform: a Junos library over NETCONF and SSH, an ORM over the device database, and Docker deployments. On it ran a health check every engineer relied on, the inventory the manager planned capacity with, a public looking glass for other providers, and the chat bot the team used to blackhole DDoS traffic (I contributed TLS support to its open-source framework). I also taught the network engineers to automate their own work.</p>

<h2 id="first-data">First Data, 2013–2015</h2>

<p class="work"><span class="t">Card authorization switch</span>, <span class="m">2013–2015.</span> C and C++ on AIX for a switch that handled up to 6,000 card authorizations a second at holiday peak, under a contractual 99.999% availability target. I built the switch-side collector for its real-time statistics monitor, brought in GCC and Clang’s static analyzer and sanitizers, and fixed thousands of lines of C without a production incident.</p>
