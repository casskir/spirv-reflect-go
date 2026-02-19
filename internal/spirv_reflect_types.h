#ifndef SPIRV_REFLECT_EXT_H
#define SPIRV_REFLECT_EXT_H

struct Scalar {
    uint32_t                        width;
    uint32_t                        signedness;
} Scalar;

struct Vector {
    uint32_t                        component_count;
} Vector;

struct Matrix {
    uint32_t                        column_count;
    uint32_t                        row_count;
    uint32_t                        stride; // Measured in bytes
} Matrix;

struct Traits {
    SpvReflectNumericTraits         numeric;
    SpvReflectImageTraits           image;
    SpvReflectArrayTraits           array;
} Traits;

#endif