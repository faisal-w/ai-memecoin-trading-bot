"""Meme generator module."""
from PIL import Image, ImageDraw, ImageFont
import os
import requests
from io import BytesIO
import config


class MemeGenerator:
    """Class for generating memes."""
    
    def __init__(self):
        """Initialize meme generator."""
        self.templates_dir = config.MEME_TEMPLATES_DIR
        self.output_dir = config.OUTPUT_DIR
        
        # Create directories if they don't exist
        os.makedirs(self.templates_dir, exist_ok=True)
        os.makedirs(self.output_dir, exist_ok=True)
    
    def create_meme(self, template_name, top_text='', bottom_text='', output_name=None):
        """
        Create a meme with text overlay.
        
        Args:
            template_name: Name of the template image file
            top_text: Text to place at the top of the image
            bottom_text: Text to place at the bottom of the image
            output_name: Optional custom output filename
            
        Returns:
            Path to the generated meme image
        """
        template_path = os.path.join(self.templates_dir, template_name)
        
        if not os.path.exists(template_path):
            raise FileNotFoundError(f"Template {template_name} not found in {self.templates_dir}")
        
        # Open template image
        img = Image.open(template_path)
        draw = ImageDraw.Draw(img)
        
        # Calculate font size based on image dimensions
        img_width, img_height = img.size
        font_size = int(img_height / 10)
        
        try:
            # Try to use a nice font
            font = ImageFont.truetype("/usr/share/fonts/truetype/dejavu/DejaVuSans-Bold.ttf", font_size)
        except:
            # Fallback to default font
            font = ImageFont.load_default()
        
        # Add top text
        if top_text:
            self._add_text(draw, top_text, img_width, font_size, 'top', font)
        
        # Add bottom text
        if bottom_text:
            self._add_text(draw, bottom_text, img_width, img_height - font_size * 2, 'bottom', font)
        
        # Generate output path
        if output_name is None:
            output_name = f"meme_{template_name}"
        output_path = os.path.join(self.output_dir, output_name)
        
        # Save the meme
        img.save(output_path)
        return output_path
    
    def _add_text(self, draw, text, x_center, y_pos, position, font):
        """Add text with outline to image."""
        # Convert text to uppercase for classic meme style
        text = text.upper()
        
        # Get text bounding box for centering
        bbox = draw.textbbox((0, 0), text, font=font)
        text_width = bbox[2] - bbox[0]
        text_height = bbox[3] - bbox[1]
        
        # Center the text
        x = (x_center - text_width) / 2
        y = y_pos
        
        # Draw text outline (black)
        outline_range = 2
        for adj_x in range(-outline_range, outline_range + 1):
            for adj_y in range(-outline_range, outline_range + 1):
                draw.text((x + adj_x, y + adj_y), text, font=font, fill='black')
        
        # Draw text (white)
        draw.text((x, y), text, font=font, fill='white')
    
    def download_template(self, url, template_name):
        """
        Download a meme template from a URL.
        
        Args:
            url: URL of the image to download
            template_name: Name to save the template as
            
        Returns:
            Path to the downloaded template
        """
        response = requests.get(url)
        if response.status_code == 200:
            template_path = os.path.join(self.templates_dir, template_name)
            with open(template_path, 'wb') as f:
                f.write(response.content)
            return template_path
        else:
            raise Exception(f"Failed to download template from {url}")
    
    def list_templates(self):
        """List all available meme templates."""
        templates = []
        if os.path.exists(self.templates_dir):
            templates = [f for f in os.listdir(self.templates_dir) 
                        if f.lower().endswith(('.png', '.jpg', '.jpeg', '.gif'))]
        return templates


def generate_meme_from_url(image_url, top_text='', bottom_text=''):
    """
    Generate a meme from an image URL.
    
    Args:
        image_url: URL of the image
        top_text: Text for top of meme
        bottom_text: Text for bottom of meme
        
    Returns:
        PIL Image object of the generated meme
    """
    response = requests.get(image_url)
    img = Image.open(BytesIO(response.content))
    
    draw = ImageDraw.Draw(img)
    img_width, img_height = img.size
    font_size = int(img_height / 10)
    
    try:
        font = ImageFont.truetype("/usr/share/fonts/truetype/dejavu/DejaVuSans-Bold.ttf", font_size)
    except:
        font = ImageFont.load_default()
    
    # Add texts
    if top_text:
        _draw_text_with_outline(draw, top_text.upper(), img_width, font_size, font)
    
    if bottom_text:
        _draw_text_with_outline(draw, bottom_text.upper(), img_width, img_height - font_size * 2, font)
    
    return img


def _draw_text_with_outline(draw, text, x_center, y_pos, font):
    """Helper function to draw text with outline."""
    bbox = draw.textbbox((0, 0), text, font=font)
    text_width = bbox[2] - bbox[0]
    
    x = (x_center - text_width) / 2
    y = y_pos
    
    # Outline
    for adj_x in range(-2, 3):
        for adj_y in range(-2, 3):
            draw.text((x + adj_x, y + adj_y), text, font=font, fill='black')
    
    # Main text
    draw.text((x, y), text, font=font, fill='white')
