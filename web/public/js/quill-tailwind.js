document.addEventListener('DOMContentLoaded', function() {
  console.log('Quill', Quill);
  const Parchment = Quill.import('parchment');

  console.log('Parchment', Parchment);

  class TailwindClass extends Parchment.Attributor {
    add(node, value) {
      if (value === 'ordered') {
        node.classList.add('list-decimal');
      } else if (value === 'bullet') {
        node.classList.add('list-disc');
      }
      return true;
    }
  }

  const TailwindListClass = new TailwindClass('list', 'ql-list', {
    scope: Parchment.Scope.BLOCK,
  });

  Quill.register(TailwindListClass, true);

  const toolbarOptions = [
    ['bold', 'italic', 'underline', 'strike'],
    ['blockquote', 'code-block'],
    [{ 'header': 1 }, { 'header': 2 }],
    [{ 'list': 'ordered' }, { 'list': 'bullet' }],
    [{ 'script': 'sub' }, { 'script': 'super' }],
    [{ 'indent': '-1' }, { 'indent': '+1' }],
    [{ 'direction': 'rtl' }],
    [{ 'size': ['small', false, 'large', 'huge'] }],
    [{ 'header': [1, 2, 3, 4, 5, 6, false] }],
    [{ 'color': [] }, { 'background': [] }],
    [{ 'font': [] }],
    [{ 'align': [] }],
    ['clean'],
    ['link']
  ];

  window.initQuill = function() {
    return new Quill('#editor', {
      modules: {
        toolbar: toolbarOptions
      },
      theme: 'snow'
    });
  };

  // Initialize Quill if the editor element exists
  const editorElement = document.getElementById('editor');
  if (editorElement) {
    const quill = initQuill();
    const form = document.querySelector('form');
    if (form) {
      form.onsubmit = function() {
        document.getElementById('template').value = quill.root.innerHTML;
      };
    }
  }
});