package lib

/*
#cgo CFLAGS: -I../../quickjs
#cgo LDFLAGS: -L../../quickjs -lquickjs
#include "quickjs.h"
#include <stdlib.h>

// 自定义的打印函数，用于console.log
static JSValue js_console_log(JSContext *ctx, JSValueConst this_val,
                              int argc, JSValueConst *argv)
{
    for (int i = 0; i < argc; i++) {
        const char *str = JS_ToCString(ctx, argv[i]);
        if (!str) {
            return JS_EXCEPTION;
        }
        if (i > 0) {
            putchar(' ');
        }
        printf("%s", str);
        JS_FreeCString(ctx, str);
    }
    putchar('\n');
    return JS_UNDEFINED;
}

// 安装console对象
static void js_std_add_helpers(JSContext *ctx)
{
    JSValue global_obj, console;
    global_obj = JS_GetGlobalObject(ctx);

    console = JS_NewObject(ctx);
    JS_SetPropertyStr(ctx, console, "log", JS_NewCFunction(ctx, js_console_log, "log", 1));
    JS_SetPropertyStr(ctx, global_obj, "console", console);

    JS_FreeValue(ctx, global_obj);
}
*/
import "C"
import (
	"fmt"
	"unsafe"
)

func JsHandler(scriptTemplate string, jsonData string) string {
	script := fmt.Sprintf(scriptTemplate, jsonData)
	//fmt.Println(script)
	// 初始化QuickJS环境
	rt := C.JS_NewRuntime()
	if rt == nil {
		fmt.Println("Could not initialize QuickJS runtime")
		return ""
	}
	defer C.JS_FreeRuntime(rt)

	// 创建上下文
	ctx := C.JS_NewContext(rt)
	if ctx == nil {
		fmt.Println("Could not create QuickJS context")
		return ""
	}
	defer C.JS_FreeContext(ctx)

	// 为上下文添加console对象
	C.js_std_add_helpers(ctx)

	// 将Go的字符串转换为C字符串
	cScript := C.CString(script)
	defer C.free(unsafe.Pointer(cScript))

	// 执行JavaScript代码
	filename := C.CString("<input>")
	defer C.free(unsafe.Pointer(filename))

	result := C.JS_Eval(ctx, cScript, C.size_t(len(script)), filename, C.JS_EVAL_TYPE_GLOBAL)
	// 获取字符串结果
	jsonString := C.JS_ToCString(ctx, result)
	return C.GoString(jsonString)
}
