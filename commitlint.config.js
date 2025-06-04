module.exports = {
    extends: ['@commitlint/config-conventional'],
    rules: {
        'type-enum': [
            2,
            'always',
            [
                'feat',     // New feature
                'fix',      // Bug fix
                'docs',     // Documentation only changes
                'style',    // Changes that don't affect code meaning
                'refactor', // Code change that neither fixes a bug nor adds a feature
                'perf',     // Performance improvement
                'test',     // Adding missing tests
                'build',    // Changes to build process
                'ci',       // Changes to CI configuration
                'chore',    // Other changes that don't modify src or test files
                'revert'    // Reverts a previous commit
            ]
        ],
        'subject-case': [2, 'never', ['upper-case', 'start-case']],
        'subject-empty': [2, 'never'],
        'subject-full-stop': [2, 'never', '.'],
        'header-max-length': [2, 'always', 72],
        'body-leading-blank': [2, 'always'],
        'footer-leading-blank': [2, 'always']
    }
};
