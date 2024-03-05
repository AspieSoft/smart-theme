onReady(async function(){
  function onInterval(){
    document.querySelectorAll('header .header-img, footer').forEach(function(elm){
      if(elm.clientHeight < 300){
        elm.classList.remove('background-size-anim');
        return;
      }

      if(elm.classList.contains('background-size-anim')){
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

              elm.classList.add('background-size-anim');
              elm.style.setProperty('--scale-offset', calculateImgSizeRatioDiff(elm, img.width, img.height) + 'px');
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
        elm.classList.remove('background-size-anim');
        return;
      }
      
      if(!elm.classList.contains('background-size-anim')){
        return;
      }
      
      const imgWidth = elm.getAttribute('img-width');
      const imgHeight = elm.getAttribute('img-height');
      if(!imgWidth || !imgHeight || imgWidth === 0 || imgHeight === 0){
        return;
      }

      elm.classList.add('background-size-anim');
      elm.style.setProperty('--scale-offset', calculateImgSizeRatioDiff(elm, imgWidth, imgHeight) + 'px');
    });
  }
  onResize();
  window.addEventListener('resize', onResize, {passive: true});


  // complex math to calculate the offset to add to a background image size
  // add this number to `background-size: auto 100%` to simulate `cover`
  function calculateImgSizeRatioDiff(elm, imgWidth, imgHeight){
    //todo: test math on bigger screens
    // console.log((imgHeight / 5))
    return (elm.clientWidth / imgWidth) * (imgWidth / imgHeight) * 100;
  }

});
