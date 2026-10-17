# Why we test inside the Dockerfile

- We wanted to verify the image other projects will use, then publish it without repeating expensive build steps.
- Our first idea was to build once, run checks inside the image, then publish that same image.
- Woodpecker can share files between steps, but cannot use an image archive directly as a step image. Running it would require a registry or a container runtime.
- We explored temporary storage and a single build/test/publish step. These still required additional setup to run the image.
- For simple tool checks, we chose Dockerfile targets instead: `artifact-test` extends the complete runtime stage and checks Go and Task through `RUN`.
- CI then builds `release`, reusing matching cached layers. Test layers stay out of the published image.
- This keeps the workflow based on standard build tooling and configuration, without our own scripts. The pattern works with Docker and Kubernetes agents; builder setup remains separate.
- We test arm64 by default; testing all published platforms is optional.
