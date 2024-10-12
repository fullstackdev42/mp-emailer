document.addEventListener('DOMContentLoaded', function() {
  const Parchment = Quill.import('parchment');

  class TailwindListClass extends Parchment.Attributor {
    add(node, value) {
      node.classList.remove('list-decimal', 'list-disc');
      if (value === 'ordered') {
        node.classList.add('list-decimal');
      } else if (value === 'bullet') {
        node.classList.add('list-disc');
      }
      return true;
    }

    remove(node) {
      node.classList.remove('list-decimal', 'list-disc');
    }
  }

  Quill.register(new TailwindListClass('list', 'ql-list', {
    scope: Parchment.Scope.BLOCK,
  }), true);

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
    // Register the htmlEditButton module
    Quill.register("modules/htmlEditButton", htmlEditButton);

    return new Quill('#editor', {
      modules: {
        toolbar: {
          container: toolbarOptions,
          handlers: {
            // Add a custom handler for the HTML edit button
            'html-edit': function() {
              const htmlEditButton = this.quill.getModule('htmlEditButton');
              if (htmlEditButton) {
                htmlEditButton.toggle();
              }
            }
          }
        },
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
    const form = document.querySelector('form');
    form?.addEventListener('submit', function(event) {
      event.preventDefault();
      document.getElementById('template').value = quill.root.innerHTML;
      this.submit();
    });
  }
});
