(function(){var r,n=document.currentScript,s=n.dataset,c=new URL(n.src),o=s.prismeUrl||c.origin,l=location.protocol,i=s.domain||location.host,a=document.referrer.replace(location.host,i),e=function(){return location.toString().replace(location.host,i)};function t(){fetch(o.concat("/api/v1/events/pageviews"),{method:"POST",headers:{"Access-Control-Max-Age":3600,"X-Prisme-Referrer":e(),"X-Prisme-Document-Referrer":a},referrerPolicy:"no-referrer-when-downgrade"}),a=e()}window.prisme={trigger:function(t,n){fetch(o.concat("/api/v1/events/custom/",t),{method:"POST",headers:{"Access-Control-Max-Age":3600,"X-Prisme-Referrer":e(),"Content-Type":"application/json"},referrerPolicy:"no-referrer-when-downgrade",body:JSON.stringify(n)})}},t(),window.history&&(r=window.history.pushState,window.history.pushState=function(){r.apply(window.history,arguments),t()},window.addEventListener("popstate",t))})()