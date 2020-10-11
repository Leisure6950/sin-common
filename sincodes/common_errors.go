/**
错误码封装
*/
package sincodes

import "github.com/sin-z/sin-common/sinerrors"

//
// 通用错误码
// 0~599
//
const (
	Success  sinerrors.Code = 0
	Redirect sinerrors.Code = 302 // 跳转

	ClientError        sinerrors.Code = 499 // 请求参数错误（调用端参数有问题）
	CSignExpired       sinerrors.Code = 451 // 签名过期
	CSignError         sinerrors.Code = 452 // 签名错误
	CSignRepeat        sinerrors.Code = 453 // 签名重复（防止重放）
	CAppKeyIsNone      sinerrors.Code = 461 // 未传app key
	CAppKeyUnsupported sinerrors.Code = 462 // 不支持的app key

	//
	// Server System Error Code
	//
	ServerError sinerrors.Code = 500 // 服务器遇到了一个未曾预料的状况，导致了它无法完成对请求的处理。一般来说，这个问题都会在服务器端的源代码出现错误时出现。

	//
	// Server net error
	//
	SNotImplemented          sinerrors.Code = 501 // 服务器不支持当前请求所需要的某个功能。当服务器无法识别请求的方法，并且无法支持其对任何资源的请求。
	SBadGateway              sinerrors.Code = 502 // 作为网关或者代理工作的服务器尝试执行请求时，从上游服务器接收到无效的响应。
	SServiceUnavailable      sinerrors.Code = 503 // 由于临时的服务器维护或者过载，服务器当前无法处理请求。这个状况是临时的，并且将在一段时间以后恢复。如果能够预计延迟时间，那么响应中可以包含一个 Retry-After 头用以标明这个延迟时间。如果没有给出这个 Retry-After 信息，那么客户端应当以处理500响应的方式处理它。
	SGatewayTimeout          sinerrors.Code = 504 // 作为网关或者代理工作的服务器尝试执行请求时，未能及时从上游服务器（URI标识出的服务器，例如HTTP、FTP、LDAP）或者辅助服务器（例如DNS）收到响应（某些代理服务器在DNS查询超时时会返回400或者500错误）。
	SHttpVersionNotSupported sinerrors.Code = 505 // 服务器不支持，或者拒绝支持在请求中使用的 HTTP 版本。这暗示着服务器不能或不愿使用与客户端相同的版本。响应中应当包含一个描述了为何版本不被支持以及服务器支持哪些协议的实体。
	SVariantAlsoNegotiates   sinerrors.Code = 506 // 由《透明内容协商协议》（RFC 2295）扩展，代表服务器存在内部配置错误：被请求的协商变元资源被配置为在透明内容协商中使用自己，因此在一个协商处理中不是一个合适的重点。
	SInsufficientStorage     sinerrors.Code = 507 // 服务器无法存储完成请求所必须的内容。这个状况被认为是临时的。WebDAV (RFC 4918)
	SBandwidthLimitExceeded  sinerrors.Code = 509 // 服务器达到带宽限制。这不是一个官方的状态码，但是仍被广泛使用。
	SNotExtended             sinerrors.Code = 510 // 获取资源所需要的策略并没有被满足。（RFC 2774）

	//
	// Server sinerrors.Code or logic exception
	//
	SException sinerrors.Code = 540 // 服务内部错误

	SAlarm sinerrors.Code = 599 // 监控系统收到此错误码，将立即报警
)

func init() {
	sinerrors.SetErrorWithToastAI(Success, "success", "操作成功")
	sinerrors.SetErrorAI(Redirect, "")
	sinerrors.SetErrorWithToastAI(ClientError, "请求参数错误", "网络异常，请重试")
	sinerrors.SetErrorWithToastAI(CSignExpired, "请求异常", "网络异常，请重试")
	sinerrors.SetErrorWithToastAI(CSignError, "请求异常", "网络异常，请重试")
	sinerrors.SetErrorWithToastAI(CSignRepeat, "请求异常", "网络异常，请重试")
	sinerrors.SetErrorWithToastAI(CAppKeyIsNone, "not found appkey", "网络异常，请重试")
	sinerrors.SetErrorWithToastAI(CAppKeyUnsupported, "unsupported appkey", "网络异常，请重试")

	sinerrors.SetErrorWithToastAI(ServerError, "内部系统错误", "系统开小差了，请稍后重试")

	sinerrors.SetErrorWithToastAI(SNotImplemented, "不支持此请求", "网络开小差了，请稍后重试")
	sinerrors.SetErrorWithToastAI(SBadGateway, "依赖接口数据错误", "网络开小差了，请稍后重试")
	sinerrors.SetErrorWithToastAI(SServiceUnavailable, "服务暂时不可用", "网络开小差了，请稍后重试")
	sinerrors.SetErrorWithToastAI(SGatewayTimeout, "访问依赖接口超时", "网络开小差了，请稍后重试")
	sinerrors.SetErrorWithToastAI(SHttpVersionNotSupported, "不支持的 HTTP 版本", "网络开小差了，请稍后重试")
	sinerrors.SetErrorWithToastAI(SVariantAlsoNegotiates, "服务内部配置错误", "网络开小差了，请稍后重试")
	sinerrors.SetErrorWithToastAI(SInsufficientStorage, "关键存储路径出错", "网络开小差了，请稍后重试")
	sinerrors.SetErrorWithToastAI(SBandwidthLimitExceeded, "服务器带宽超限", "网络开小差了，请稍后重试")
	sinerrors.SetErrorWithToastAI(SNotExtended, "资源依赖访问策略不足", "网络开小差了，请稍后重试")

	sinerrors.SetErrorWithToastAI(SException, "服务内部错误", "系统开小差了，工程师们正在紧急修复，请稍后重试")

	sinerrors.SetErrorWithToastAI(SAlarm, "服务出现严重故障", "数据异常，请联系客服")
}
