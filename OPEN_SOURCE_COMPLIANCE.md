# Open Source Compliance

Author: 武汉大学开源软件与技术课程 2026

## Project License

urlAPI is distributed under GPL-3.0. This is a strong copyleft license. If a
modified version or combined derivative work is distributed, the distributor
must provide the corresponding source code under GPL-compatible terms and keep
the license and copyright notices.

GPL-3.0 does not contain the AGPL network-service clause. Pure internal SaaS
operation does not trigger the same source-disclosure obligation as AGPL, but
binary, Docker image, appliance, or other software distribution does trigger
GPL-3.0 compliance duties.

## Current Compliance Baseline

- `LICENSE` is present and contains GPL-3.0.
- `NOTICE` is present for project-level notices and release reminders.
- Backend dependency inventory is declared in `go.mod` and `go.sum`.
- Frontend dependency inventory is declared in `static/package.json` and
  `static/package-lock.json`.
- The project integrates upstream model and content services; their API terms,
  model licenses, and usage policies must be reviewed separately.

## Required Release Checklist

1. Generate an SBOM from the exact release source tree and dependency lockfiles.
2. Run SCA/license scanning on Go modules, npm packages, vendored files, static
   assets, fonts, images, and copied snippets.
3. Confirm no GPL/AGPL-incompatible dependency is introduced for the intended
   distribution model.
4. Bundle `LICENSE`, `NOTICE`, and generated third-party notices with source,
   Docker images, binaries, and frontend bundles.
5. If files from third-party projects are modified, keep visible modification
   notices in those files or in the release notice set.
6. Keep Git history, release tags, build logs, dependency lockfiles, and scan
   reports as evidence for later audits.

## License Admission Policy

- Low risk: MIT, BSD, Apache-2.0. Keep notices and disclaimers.
- Medium risk: LGPL, MPL. Review linking and modified-component obligations.
- High risk: GPL, AGPL. Legal review is required before combining with code
  intended for closed-source commercial distribution; AGPL also affects network
  service deployment.

## AI And Data Compliance

- Check provider terms for OpenAI-compatible, Anthropic, Gemini, Azure,
  Moonshot, Alibaba, DeepSeek, and other configured upstream services.
- Check model-community licenses for MAU limits, field-of-use restrictions,
  anti-competitive-model clauses, and attribution requirements.
- Do not assume a model described as "open" is OSI open source.
- Keep prompts, generated assets, training/fine-tuning data, and source URLs
  traceable when they are reused in products or datasets.
- Avoid sending private, personal, or confidential data to upstream providers
  unless the deployment policy and user consent allow it.
