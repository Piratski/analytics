(function(){var f,m=document.currentScript,e=location,t=m.dataset,C=new URL(m.src),a="no-referrer-when-downgrade",r="POST",v=e.protocol,l=t.prismeUrl||C.origin,d=t.domain||e.host,x=t.path||e.pathname,o=!!t.manual&&t.manual!=="false"||!1,w=t.visitorId,h=t.outboundLinks!=="false",y=t.fileDownloads!=="false",j=(t.extraDownloadsFileTypes||"").split(","),p=document.referrer.replace(e.host,d),g=1,n=globalThis,b="Request"in n&&"keepalive"in new Request(""),E=["7z","avi","csv","dmg","docx","exe","gz","key","midi","mov","mp3","mp4","mpeg","pdf","pkg","pps","ppt","pptx","rar","rtf","txt","wav","wma","wmv","xlsx","zip"].concat(j);function i(t){return t||(t={}),t.domain||(o?t.domain=e.host:t.domain=d),t.path||(o||g>1?t.path=e.pathname:t.path=x),t.visitorId||(t.visitorId=w),t.url=v.concat("//",t.domain,t.path,e.search),t}function c(e,t){return t["Access-Control-Max-Age"]=3600,t["X-Prisme-Referrer"]=e.url,e.visitorId&&(t["X-Prisme-Visitor-Id"]=e.visitorId.toString()),t}function _(e,t){if(e.defaultPrevented)return!1;var n=!t.target||t.target.match(/^_(self|parent|top)$/i),s=!(e.ctrlKey||e.metaKey||e.shiftKey)&&e.type==="click";return n&&s}function s(e){e=i(e),fetch(l.concat("/api/v1/events/pageviews"),{method:r,headers:c(e,{"X-Prisme-Document-Referrer":p}),keepalive:!0,referrerPolicy:a}),p=e.url,g++}function O(e,t){return t=i(t),fetch(l.concat("/api/v1/events/clicks/outbound-link"),{method:r,headers:c(t,{}),keepalive:!0,referrerPolicy:a,body:e})}function u(t){if(t.type==="auxclick"&&t.button!==1||!(t.target instanceof Element))return;var s,o,i,a=t.target.closest("a");if(!a)return;if(s=new URL(a.href||"",e.origin),s.search="",h&&s.host!==e.host){o=!b&&_(t,s),i=!1;function r(){!i&&o&&(i=!0,n.location.assign(s))}console.log("follow link manually",o),o&&(t.preventDefault(),setTimeout(r,5e3)),O(s).finally(r)}}(h||y)&&(document.addEventListener("click",u),document.addEventListener("auxclick",u)),n.prisme={pageview:s,trigger(e,t,n){n=i(n),fetch(l.concat("/api/v1/events/custom/",e),{method:r,headers:c(n,{"Content-Type":"application/json"}),keepalive:!0,referrerPolicy:a,body:JSON.stringify(t)})}},o||(delete n.prisme.pageview,s(),n.history&&(f=n.history.pushState,n.history.pushState=function(){f.apply(n.history,arguments),s()},n.addEventListener("popstate",s)))})()