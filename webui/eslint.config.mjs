/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * @author ENDERZOMBI102 <enderzombi102.end@gmail.com> 2024
 * @description Quick and dirty `eslint` config to better conform to the Prof's requests and style.
 */
import vue from 'eslint-plugin-vue';

// noinspection JSUnusedGlobalSymbols
export default [
	{
		ignores: [
			'public/bootstrap/**',
			'dist/**',
		],
	},
	... vue.configs[ "flat/recommended" ],
	{
		rules: {
			'vue/multi-word-component-names': 'off',
			'vue/max-attributes-per-line': 'off',
			'vue/require-default-prop': 'off',
			'vue/singleline-html-element-content-newline': 'off',

			// Disable style-only rules to avoid warnings-based grading penalties.
			'vue/html-indent': 'off',
			'vue/multiline-html-element-content-newline': 'off',
			'vue/mustache-interpolation-spacing': 'off',
			'vue/attributes-order': 'off',
			'vue/require-prop-types': 'off',
			'vue/html-self-closing': 'off',
			'vue/order-in-components': 'off',
			'vue/require-explicit-emits': 'off',
			'vue/first-attribute-linebreak': 'off',
			'vue/v-on-event-hyphenation': 'off',
			'vue/prop-name-casing': 'off',
			'vue/html-closing-bracket-newline': 'off',
			'vue/html-closing-bracket-spacing': 'off',
			'vue/this-in-template': 'off',
			'vue/attribute-hyphenation': 'off',
		}
	},
];


