(function(){var a,c=document.currentScript,e=location,t=c.dataset,v=new URL(c.src),d="no-referrer-when-downgrade",h="POST",b=e.protocol,r=t.prismeUrl||v.origin,o=t.domain||e.host,g=t.path||e.pathname,s=!!t.manual||!1,p=t.visitorId,f=t.visitorAnonymous,m=document.referrer.replace(e.host,o),i=1;function u(t){return t||(t={}),t.domain||(s?t.domain=e.host:t.domain=o),t.path||(s||i>1?t.path=e.pathname:t.path=g),t.visitorId||(t.visitorId=p),"anonymous"in t||(t.anonymous=f),t.url=b.concat("//",t.domain,t.path,e.search),t}function l(e,t){return t["Access-Control-Max-Age"]=3600,t["X-Prisme-Referrer"]=e.url,e.visitorId&&(t["X-Prisme-Visitor-Id"]=e.visitorId.toString(),e.anonymous==!0&&(t["X-Prisme-Visitor-Anon"]="1")),t}function n(e){e=u(e),fetch(r.concat("/api/v1/events/pageviews"),{method:h,headers:l(e,{"X-Prisme-Document-Referrer":m}),keepalive:!0,referrerPolicy:d}),m=e.url,i++}window.prisme={pageview:n,trigger(e,t,n){n=u(n),fetch(r.concat("/api/v1/events/custom/",e),{method:h,headers:l(n,{"Content-Type":"application/json"}),keepalive:!0,referrerPolicy:d,body:JSON.stringify(t)})}},s||(delete window.prisme.pageview,n(),window.history&&(a=window.history.pushState,window.history.pushState=function(){a.apply(window.history,arguments),n()},window.addEventListener("popstate",n)))})()