const OFF = 0;
// const WARN = 1;
const ERROR = 2;

module.exports = {
  extends: ["eslint:recommended"],
  parserOptions: {
    ecmaVersion: 2020,
  },
  env: {
    es6: true,
    node: true,
  },
  rules: {
    indent: OFF,
    semi: ERROR,
    quotes: [ERROR, "double", {avoidEscape: true, allowTemplateLiterals: true}],
    "no-empty": ERROR,
    "array-callback-return": ERROR,
    "consistent-return": ERROR,
    eqeqeq: OFF,
    "prefer-const": ERROR,
    "no-unused-vars": [ERROR, {args: "none", varsIgnorePattern: "^_"}],
    "no-console": OFF,
    "no-debugger": OFF,
    "require-atomic-updates": OFF,
  },
  globals: {
  },
};
