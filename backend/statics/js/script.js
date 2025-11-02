document.addEventListener('DOMContentLoaded', () => {
  const boton = document.getElementById('miBoton');
  const mensaje = document.getElementById('mensaje');

  boton.addEventListener('click', () => {
    mensaje.style.display = 'block';
  });
});
