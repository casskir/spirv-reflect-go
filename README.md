## SPIR-V Reflect Go (CGO wrapper)

This repository provides Go bindings for [SPIRV-Reflect](https://github.com/KhronosGroup/SPIRV-Reflect), a lightweight C library for reflecting Vulkan SPIR-V shader modules. It lets you inspect shader entry points, descriptor sets/bindings, push constants, input/output variables, and other metadata directly from compiled `.spv` bytecode.

### Highlights:
- Thin CGO bridge over the upstream SPIRV-Reflect C API
- Simple, idiomatic Go types exposed in `pkg/spvreflect`
- Example CLI (`cmd/spirv-info`) demonstrating basic reflection

### Requirements
- Go with CGO enabled (a working C toolchain is required)
- Git submodules (to fetch the embedded SPIRV-Reflect sources)

Initialize submodules:
```
git submodule update --init --recursive
```

### Installation
This module’s path is currently `spvreflect` (see `go.mod`). If you are using it inside this repository, you can import it as shown in the examples below. If you vendored or forked the project under your own module path, replace the import path accordingly.

Inside this repository (local module):
```
import "spvreflect/pkg/spvreflect"
```

### Quick start (code example)
The following snippet shows how to load a `.spv` file and print some basic metadata. It mirrors the sample in `cmd/spirv-info`.

```
package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "os"
    "unsafe"

    "spvreflect/pkg/spvreflect"
)

func main() {
    if len(os.Args) != 2 {
        log.Fatalf("Usage: %s <shader.spv>", os.Args[0])
    }

    // 1) Read SPIR-V bytecode
    bytecode, err := ioutil.ReadFile(os.Args[1])
    if err != nil {
        log.Fatalf("failed to read shader file: %v", err)
    }

    // 2) Create a ShaderModule from the bytecode
    var module spvreflect.ShaderModule
    result := spvreflect.CreateShaderModule(
        uint(len(bytecode)),
        unsafe.Pointer(&bytecode[0]),
        &module,
    )
    if result != spvreflect.ResultSuccess {
        log.Fatalf("failed to reflect shader: error code %d", result)
    }
    defer spvreflect.DestroyShaderModule(&module) // free C allocations

    // Convert internal C data to Go fields where applicable
    module.Deref()

    // 3) Print basic information
    fmt.Printf("Entry Point: %s\n", module.EntryPointName)
    fmt.Printf("Source Language: %v\n", spvreflect.SourceLanguageName(module.SourceLanguage))
    fmt.Printf("Shader Stage: %v\n", module.ShaderStage)

    // 4) Enumerate descriptor bindings
    descriptorBindings := make([]*spvreflect.DescriptorBinding, module.DescriptorBindingCount)
    spvreflect.EnumerateDescriptorBindings(&module, &module.DescriptorBindingCount, descriptorBindings)

    fmt.Println("\nDescriptor Bindings:")
    for _, b := range descriptorBindings {
        b.Deref() // load Go-side copies of reflected fields
        fmt.Printf("\tSet: %d, Binding: %d, Count: %d\n", b.Set, b.Binding, b.Count)
        fmt.Printf("\tType: %v, Name: %s\n", b.DescriptorType, b.Name)
    }
}
```

### Build and run the example CLI
An example utility is available at `cmd/spirv-info`. It prints basic information about a SPIR-V module.

Build:
```
go build -o bin/spirv-info ./cmd/spirv-info
```

Run against the included sample shader:
```
./bin/spirv-info ./cmd/spirv-info/sprite.vert.spv
```

### API notes
- Memory management: many functions allocate memory inside the C library. Always call `spvreflect.DestroyShaderModule(&module)` once you are done. Some wrapper types provide a `Deref()` method to copy/convert fields into Go-managed memory; call it before reading their fields.
- Error handling: most creation/enumeration functions return a `Result` code (e.g., `ResultSuccess`). Check these results.

### Troubleshooting
- CGO disabled: ensure `CGO_ENABLED=1` (the default on most platforms with a C toolchain installed).
- Missing submodule: run `git submodule update --init --recursive` before building.
- Compiler/linker errors: install a C compiler toolchain (e.g., `build-essential` on Debian/Ubuntu, Xcode CLT on macOS, or MSYS2/MinGW on Windows).

### Code Generation
The Go bindings are generated using [c-for-go](https://github.com/xlab/c-for-go) tool. To regenerate the bindings:
```bash
c-for-go -ccdefs -out pkg/ spirv-reflect.yml
```

### License and acknowledgements
- This project wraps the upstream SPIRV-Reflect library (see `internal/SPIRV-Reflect`). Please refer to its license for upstream terms.
- SPIR-V and Vulkan are trademarks of the Khronos Group Inc.
