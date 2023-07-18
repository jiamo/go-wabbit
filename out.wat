(module
(import "env" "_printi" (func $_printi ( param i32 )))
(import "env" "_printf" (func $_printf ( param f64 )))
(import "env" "_printb" (func $_printb ( param i32 )))
(import "env" "_printc" (func $_printc ( param i32 )))
(global $a (mut i32) (i32.const 0))
(global $b (mut i32) (i32.const 0))
(global $c (mut i32) (i32.const 0))
(func $main (export "main")


block $return

i32.const 1
global.set $a
i32.const 2
global.set $b
global.get $a
global.get $b
i32.add
i32.const 3
i32.add
global.set $c
global.get $c
drop
global.get $c
call $_printi
global.get $a
global.get $b
i32.add
call $_printi
global.get $a
global.get $b
i32.sub
call $_printi
global.get $b
global.get $c
i32.mul
call $_printi
global.get $c
global.get $b
i32.div_s
call $_printi
global.get $a
call $_printi
i32.const 0
global.get $a
i32.sub
call $_printi
i32.const 0
global.get $a
i32.sub
global.get $b
i32.add
call $_printi
global.get $a
global.get $b
global.get $c
i32.mul
i32.add
call $_printi
i32.const 42
global.set $c
global.get $c
i32.const 5
i32.add
call $_printi
global.get $c
call $_printi
end
)

)
