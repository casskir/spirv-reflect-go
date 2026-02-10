package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"unsafe"

	"github.com/casskir/spirv-reflect-go/pkg/spvreflect"
)

func main() {
	// 1. Read SPIR-V file into byte slice
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s <shader.spv>", os.Args[0])
	}

	bytecode, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalf("failed to read shader file: %v", err)
	}

	// 2. Create an instance of ShaderModule structure
	// Depending on your renaming rules, it might be called ShaderModule
	var module spvreflect.ShaderModule

	// 3. Call module creation function
	// spvReflectCreateShaderModule(size, data, &module)
	result := spvreflect.CreateShaderModule(
		uint(len(bytecode)),
		unsafe.Pointer(&bytecode[0]), // передаем указатель на первый байт
		&module,
	)
	// Check result (SUCCESS constant has also gone through your renaming rules)
	if result != spvreflect.ResultSuccess {
		log.Fatalf("failed to reflect shader: error code %d", result)
	}

	module.Deref()

	// Don't forget to free C memory at the end
	defer spvreflect.DestroyShaderModule(&module)

	// 4. Output basic information
	fmt.Printf("Entry Point: %s\n", module.EntryPointName)
	fmt.Printf("Source Language: %v\n", spvreflect.SourceLanguageName(module.SourceLanguage))
	fmt.Printf("Shader Stage: %v\n", module.ShaderStage)

	// 5. Iterate through descriptors (Binding)
	descriptorBindings := make([]*spvreflect.DescriptorBinding, module.DescriptorBindingCount)
	spvreflect.EnumerateDescriptorBindings(&module, &module.DescriptorBindingCount, descriptorBindings)

	fmt.Println("\nDescriptor Bindings:")
	for _, b := range descriptorBindings {
		b.Deref()

		fmt.Printf("\tSet: %d, Binding: %d, Count: %d\n", b.Set, b.Binding, b.Count)
		fmt.Printf("\tType: %v, Name: %s\n", b.DescriptorType, b.Name)
	}
}
