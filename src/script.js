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
  const body = document.body || document.querySelector('body');
  const headerTop = document.querySelector('header .header-top');
  const headerTopNav = document.querySelector('header .header-top nav.top-nav > ul');

  if(body && window.location.pathname === '/'){
    body.classList.add('home');
  }


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


  function onResize(){
    if(headerTop && headerTopNav){
      let headerRect = headerTop.getBoundingClientRect();

      if(headerTop.scrollWidth > headerRect.width){
        let allHidden = true;

        for(let i = headerTopNav.children.length - 1; i >= 0; i--){
          headerTopNav.children[i].style['display'] = 'none';
          headerRect = headerTop.getBoundingClientRect();
          if(headerTop.scrollWidth <= headerRect.width){
            allHidden = false;
            break;
          }
        }

        if(allHidden){
          headerTop.querySelectorAll('nav:not(.top-nav) > ul').forEach(function(elm){
            if(!allHidden){
              return;
            }

            if(headerTop.scrollWidth > headerRect.width){
              for(let i = 0; i < elm.children.length; i++){
                elm.children[i].style['display'] = 'none';
                headerRect = headerTop.getBoundingClientRect();
                if(headerTop.scrollWidth <= headerRect.width){
                  allHidden = false;
                  break;
                }
              }
            }
          });
        }
      }else{
        let allVisible = true;

        headerTop.querySelectorAll('nav:not(.top-nav) > ul').forEach(function(elm){
          if(!allVisible){
            return;
          }

          if(headerTop.scrollWidth <= headerRect.width){
            for(let i = elm.children.length - 1; i >= 0; i--){
              elm.children[i].style['display'] = '';
              headerRect = headerTop.getBoundingClientRect();
              if(headerTop.scrollWidth > headerRect.width){
                elm.children[i].style['display'] = 'none';
                allVisible = false;
                break;
              }
            }
          }
        });

        if(allVisible){
          for(let i = 0; i < headerTopNav.children.length; i++){
            headerTopNav.children[i].style['display'] = '';
            headerRect = headerTop.getBoundingClientRect();
            if(headerTop.scrollWidth > headerRect.width){
              headerTopNav.children[i].style['display'] = 'none';
              break;
            }
          }
        }
      }
    }
  }
  onResize();
  window.addEventListener('resize', onResize, {passive: true});


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
