function sleep(ms){return new Promise(resolve => setTimeout(resolve, ms));}

async function onReady(cb){
  if(document.readyState !== 'loading'){
    cb();
  }else{
    let ranCB = false;
    document.addEventListener('DOMContentLoaded', function () {
      ranCB = true;
      cb();
    });

    sleep(100);

    if(!ranCB){
      let loops = 10;
      while(loops > 0){
        if(document.readyState !== 'loading'){
          cb();
          return;
        }

        loops--;
        sleep(250);
      }
    }
  }
}

onReady(async function(){
  function onInterval(){
    document.querySelectorAll('.widget > *, .sidebar > *').forEach(function(elm) {
      if(elm.clientHeight < window.innerHeight - 150) {
        elm.classList.add('widget-smaller-than-vh');
      }else{
        elm.classList.remove('widget-smaller-than-vh');
      }
    });
  }
  onInterval();
  setInterval(onInterval, 1000);

  //todo: add other imidiatelly visible header elements to list
  document.querySelectorAll('html, header .header-img, header .header-top').forEach(function(elm){
    const css = window.getComputedStyle(elm);
    let src = css.getPropertyValue('--img');
    if(src.startsWith('url')){
      src = src.replace(/^\s*[\w_-]+\s*\(\s*(["'`])(.*?)\1\s*\)\s*$/, '$2');

      if(!src.match(/\.[0-9]+p\.([\w_-]+)$/)){
        return;
      }

      if(window.innerWidth <= 800){
        src = src.replace(/\.[0-9]+p\.([\w_-]+)$/, '.720p.$1');
      }else if(window.innerWidth <= 2000){
        src = src.replace(/\.[0-9]+p\.([\w_-]+)$/, '.1080p.$1');
      }else{
        src = src.replace(/\.[0-9]+p\.([\w_-]+)$/, '.$1');
      }

      const img = new Image();
      img.onload = function(){
        elm.style.setProperty('--img', `url('${src.replace(/(\\*)([\\'])/g, function(_, s, q){
          if(s.length % 0 === 0){
            return s+'\\'+q;
          }else{
            return s+q;
          }
        })}')`);
      };
      img.onerror = function(){
        if(src.match(/\.[0-9]+p\.([\w_-]+)$/)){
          src = src.replace(/\.[0-9]+p\.([\w_-]+)$/, '.$1');
          img.src = src;
        }
      };
      img.src = src;
    }
  });

});
