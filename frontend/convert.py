import os
import re

def html_to_jsx(html):
    # Basic HTML to JSX conversions
    # Replace class= with className=
    html = re.sub(r'\bclass=', 'className=', html)
    # Replace for= with htmlFor=
    html = re.sub(r'\bfor=', 'htmlFor=', html)
    
    # SVG properties
    svg_props = [
        'stroke-width', 'stroke-linecap', 'stroke-linejoin', 'fill-rule', 'clip-rule', 
        'stroke-miterlimit', 'stop-color', 'stop-opacity'
    ]
    for prop in svg_props:
        camel_prop = ''.join(word.capitalize() if i > 0 else word for i, word in enumerate(prop.split('-')))
        html = re.sub(r'\b' + prop + r'=', camel_prop + '=', html)
        
    # Close self-closing tags
    html = re.sub(r'<img([^>]+?)(?<!/)>', r'<img\1 />', html)
    html = re.sub(r'<input([^>]+?)(?<!/)>', r'<input\1 />', html)
    html = re.sub(r'<br([^>]+?)(?<!/)>', r'<br\1 />', html)
    html = re.sub(r'<hr([^>]+?)(?<!/)>', r'<hr\1 />', html)
    
    # Remove HTML comments to avoid JSX parse errors
    html = re.sub(r'<!--(.*?)-->', '', html, flags=re.DOTALL)
    
    # Inline styles: style="transform: translateY(-1px);" -> style={{ transform: 'translateY(-1px)' }}
    # For simplicity, we just strip inline styles if they are too complex, but Stitch rarely uses them except for SVG.
    # Actually Stitch HTML often doesn't have complex inline styles, they use Tailwind.
    html = re.sub(r'style="([^"]*)"', 'style={{}}', html) # Just clear inline styles for now
    
    # Extract body content
    body_match = re.search(r'<body[^>]*>(.*)</body>', html, re.DOTALL | re.IGNORECASE)
    if body_match:
        content = body_match.group(1)
        # Remove <script> tags
        content = re.sub(r'<script.*?>.*?</script>', '', content, flags=re.DOTALL | re.IGNORECASE)
        return content.strip()
    return html

def create_component(name, jsx_content):
    return f"""import React from 'react';
import {{ Link }} from 'react-router-dom';

export default function {name}() {{
  return (
    <>
      {jsx_content}
    </>
  );
}}
"""

files = {
    'home.html': ('HomePage', 'HomePage'),
    'admin_dashboard.html': ('AdminDashboard', 'AdminPage'),
    'admin_sources.html': ('AdminSources', 'AdminPage'),
    'admin_content.html': ('AdminContent', 'AdminPage'),
    'article_detail.html': ('ArticlePage', 'ArticlePage'),
    'login.html': ('LoginPage', 'LoginPage'),
}

in_dir = '../.stitch/designs'
out_dir = 'src/pages'

for filename, (comp_name, folder) in files.items():
    filepath = os.path.join(in_dir, filename)
    if os.path.exists(filepath):
        with open(filepath, 'r', encoding='utf-8') as f:
            html = f.read()
            jsx_content = html_to_jsx(html)
            
            # Wrap in component
            final_code = create_component(comp_name, jsx_content)
            
            out_folder = os.path.join(out_dir, folder)
            os.makedirs(out_folder, exist_ok=True)
            out_file = os.path.join(out_folder, f'{comp_name}.jsx')
            
            with open(out_file, 'w', encoding='utf-8') as out_f:
                out_f.write(final_code)
            print(f"Generated {comp_name}.jsx")
    else:
        print(f"File not found: {filepath}")
