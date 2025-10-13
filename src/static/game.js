// Simple JS pour gérer la sélection de pion et stockage local
document.addEventListener('DOMContentLoaded', function () {
  const buttons = document.querySelectorAll('.pion');
  const selectedSpan = document.getElementById('selected-player');

  function setSelected(value, label) {
    localStorage.setItem('selectedPion', value);
    selectedSpan.textContent = label;
  }

  // charger valeur sauvegardée
  const saved = localStorage.getItem('selectedPion');
  if (saved) {
    const btn = document.querySelector('.pion[data-value="' + saved + '"]');
    if (btn) selectedSpan.textContent = btn.textContent;
  }

  buttons.forEach(btn => {
    btn.addEventListener('click', () => {
      const val = btn.getAttribute('data-value');
      setSelected(val, btn.textContent);
      // Visuel
      buttons.forEach(b => b.classList.remove('active'));
      btn.classList.add('active');
    });
  });
});
