onReady(async function(){
  function onInterval(){
    document.querySelectorAll('header .header-img, footer').forEach(function(elm){
      if(elm.clientHeight < 300){
        elm.classList.remove('background-size-w', 'background-size-h');
        return;
      }

      if(elm.classList.contains('background-size-w') || elm.classList.contains('background-size-h')){
        return;
      }

      const css = window.getComputedStyle(elm);
      if(typeof css.getPropertyValue('--img-size') !== 'string' || typeof css.getPropertyValue('--img') !== 'string'){
        return;
      }

      if(css.getPropertyValue('--img-size').includes('cover')){
        css.getPropertyValue('--img').replace(/url\s*\(\s*(["'`])(.*?)\1\s*\)/i, function(_, _, url){
          const img = new Image();
          img.src = url;
          img.onload = function(){
            if(img.width !== 0 && img.height !== 0){
              elm.setAttribute('img-width', img.width);
              elm.setAttribute('img-height', img.height);
              if(elm.clientWidth / img.width >= elm.clientHeight / img.height){
                elm.classList.add('background-size-w');
                elm.classList.remove('background-size-h');
              }else{
                elm.classList.add('background-size-h');
                elm.classList.remove('background-size-w');
              }
            }
            img.remove();
          }
        });
      }
    });

    document.querySelectorAll('header .header-top').forEach(function(elm) {
      const css = window.getComputedStyle(elm);
      const img = css.getPropertyValue('--img');
      if(img !== '' && img !== 'none' && img.match(/^\s*([\w_-]+-gradient|url)/)){
        elm.style['background-color'] = 'transparent';
        elm.style['animation-name'] = 'scroll-header-top-shadow';
      }else{
        elm.style['background-color'] = '';
        elm.style['animation-name'] = '';
      }
    });
  }
  onInterval();
  setInterval(onInterval, 1000);


  function onResize(){
    document.querySelectorAll('header .header-img, footer').forEach(function(elm){
      if(elm.clientHeight < 300){
        elm.classList.remove('background-size-w', 'background-size-h');
        return;
      }

      if(!elm.classList.contains('background-size-w') && !elm.classList.contains('background-size-h')){
        return;
      }

      const imgWidth = elm.getAttribute('img-width');
      const imgHeight = elm.getAttribute('img-height');
      if(!imgWidth || !imgHeight || imgWidth === 0 || imgHeight === 0){
        return;
      }

      if(imgWidth >= imgHeight){
        //todo: fix background-size-w and background-size-h not setting correctly
      }

      // console.log(elm.clientWidth / imgWidth, elm.clientHeight / imgHeight)
      if(elm.clientWidth / imgWidth >= elm.clientHeight / imgHeight){
        elm.classList.add('background-size-w');
        elm.classList.remove('background-size-h');
      }else{
        elm.classList.add('background-size-h');
        elm.classList.remove('background-size-w');
      }
    });
  }
  onResize();
  window.addEventListener('resize', onResize, {passive: true});

});
