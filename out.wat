(module
(import "env" "_printi" (func $_printi ( param i32 )))
(import "env" "_printf" (func $_printf ( param f64 )))
(import "env" "_printb" (func $_printb ( param i32 )))
(import "env" "_printc" (func $_printc ( param i32 )))
(global $a (mut i32) (i32.const 0))
(global $b (mut i32) (i32.const 0))
(func $main (export "main")


block $return

i32.const 1
global.set $a
i32.const 0
global.set $b
global.get $a
call $_printb
global.get $b
call $_printb
global.get $a
global.get $b
i32.eq
call $_printb
global.get $a
global.get $b
i32.ne
call $_printb
block $begin.0 (result i32)
block $and_block.1
global.get $a
i32.const 1
i32.xor
br_if $and_block.1
global.get $b
br $begin.0
end
i32.const 0
br $begin.0
end
call $_printb
block $begin.2 (result i32)
block $or_block.3
global.get $a
br_if $or_block.3
global.get $b
br $begin.2
end
i32.const 1
br $begin.2
end
call $_printb
global.get $a
i32.const 1
i32.xor
call $_printb
global.get $b
i32.const 1
i32.xor
call $_printb
end
)

)
