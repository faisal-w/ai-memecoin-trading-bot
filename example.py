"""Example usage of the meme generator."""
from meme_generator import MemeGenerator
import os

def main():
    """Example usage of meme generator."""
    # Initialize the generator
    meme_gen = MemeGenerator()
    
    print("Meme Generator Example")
    print("=" * 50)
    
    # List available templates
    templates = meme_gen.list_templates()
    print(f"\nAvailable templates: {len(templates)}")
    for template in templates:
        print(f"  - {template}")
    
    # Example: Create a meme from a template (if templates exist)
    if templates:
        template_name = templates[0]
        print(f"\nGenerating meme from template: {template_name}")
        
        try:
            output_path = meme_gen.create_meme(
                template_name=template_name,
                top_text="This is the top text",
                bottom_text="This is the bottom text"
            )
            print(f"✅ Meme generated successfully: {output_path}")
        except Exception as e:
            print(f"❌ Error generating meme: {e}")
    else:
        print("\n⚠️  No templates found!")
        print("Add image files to the 'meme_templates' directory to get started.")
    
    # Example: Download a template from URL
    print("\nYou can also download templates from URLs:")
    print("Example code:")
    print("""
    meme_gen.download_template(
        url='https://example.com/meme-template.jpg',
        template_name='my_template.jpg'
    )
    """)
    
    print("\n" + "=" * 50)
    print("Check the 'output' directory for generated memes!")


if __name__ == '__main__':
    main()
