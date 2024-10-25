document.addEventListener('DOMContentLoaded', function() {
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
    ['link'],
  ];

  window.initQuill = function() {
    const editorElement = document.getElementById('editor');
    if (!editorElement) {
      console.error('Editor element not found');
      return null;
    }
    return new Quill('#editor', {
      modules: {
        toolbar: toolbarOptions,
        htmlEditButton: {
          buttonHTML: "&lt;&gt;",
          buttonTitle: "Edit HTML",
        }
      },
      theme: 'snow'
    });
  };

  const editorElement = document.getElementById('editor');
  if (editorElement) {
    const quill = initQuill();
    if (quill) {
      const form = document.querySelector('form');
      form?.addEventListener('submit', function(event) {
        event.preventDefault();
        const templateElement = document.getElementById('template');
        if (templateElement) {
          templateElement.value = quill.root.innerHTML;
          this.submit();
        } else {
          console.error('Template element not found');
        }
      });
    }
  } else {
    console.warn('Editor element not found, Quill initialization skipped');
  }
});
