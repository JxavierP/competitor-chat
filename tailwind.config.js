// @type {import('tailwindcss').Config} */

export const content = ["./web/components/**/*.templ"];
export const theme = {
  extend: {
    typography: () => ({
      pine: {
        css: {
          '--tw-prose-body': 'var(--color-rose-pine-text)',
          '--tw-prose-headings': 'var(--color-rose-pine-text)',
          '--tw-prose-lead': 'var(--color-rose-pine-text)',
          '--tw-prose-links': 'var(--color-rose-pine-iris)',
          '--tw-prose-bold': 'var(--color-rose-pine-text)',
          '--tw-prose-counters': 'var(--color-rose-pine-text)',
          '--tw-prose-bullets': 'var(--color-rose-pine-text)',
          '--tw-prose-hr': 'var(--color-rose-highlight-high)',
          '--tw-prose-quotes': 'var(--color-rose-pine-subtle)',
          '--tw-prose-quote-borders': 'var(--color-rose-pine-highlight-med)',
          '--tw-prose-captions': 'var(--color-rose-pine-text)',
          '--tw-prose-code': 'var(--color-rose-pine-text)',
          '--tw-prose-pre-code': 'var(--color-rose-pine-iris)',
          '--tw-prose-pre-bg': 'var(--color-rose-highlight-low)',
          '--tw-prose-th-borders': 'var(--color-rose-highlight-high)',
          '--tw-prose-td-borders': 'var(--color-rose-highlight-med)',
        },
      },
    }),
  },
};

