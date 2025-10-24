use yew::prelude::*;
use web_sys::HtmlInputElement;
use wasm_bindgen_futures::spawn_local;
use gloo_net::http::Request;
use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize)]
struct MemeRequest {
    template: String,
    image_url: String,
    top_text: String,
    bottom_text: String,
}

#[derive(Debug, Deserialize)]
struct TemplatesResponse {
    templates: Vec<String>,
}

enum Msg {
    UpdateTopText(String),
    UpdateBottomText(String),
    UpdateImageUrl(String),
    SelectTemplate(String),
    GenerateMeme,
    TemplatesLoaded(Vec<String>),
    MemeGenerated(String),
    Error(String),
}

struct App {
    top_text: String,
    bottom_text: String,
    image_url: String,
    selected_template: String,
    templates: Vec<String>,
    generated_meme_url: Option<String>,
    loading: bool,
    error_message: Option<String>,
}

impl Component for App {
    type Message = Msg;
    type Properties = ();

    fn create(ctx: &Context<Self>) -> Self {
        // Load templates on startup
        let link = ctx.link().clone();
        spawn_local(async move {
            match Request::get("/api/templates").send().await {
                Ok(response) => {
                    if let Ok(data) = response.json::<TemplatesResponse>().await {
                        link.send_message(Msg::TemplatesLoaded(data.templates));
                    }
                }
                Err(_) => {
                    link.send_message(Msg::Error("Failed to load templates".to_string()));
                }
            }
        });

        Self {
            top_text: String::new(),
            bottom_text: String::new(),
            image_url: String::new(),
            selected_template: String::new(),
            templates: Vec::new(),
            generated_meme_url: None,
            loading: false,
            error_message: None,
        }
    }

    fn update(&mut self, ctx: &Context<Self>, msg: Self::Message) -> bool {
        match msg {
            Msg::UpdateTopText(text) => {
                self.top_text = text;
                true
            }
            Msg::UpdateBottomText(text) => {
                self.bottom_text = text;
                true
            }
            Msg::UpdateImageUrl(url) => {
                self.image_url = url;
                self.selected_template.clear();
                true
            }
            Msg::SelectTemplate(template) => {
                self.selected_template = template;
                self.image_url.clear();
                true
            }
            Msg::GenerateMeme => {
                self.loading = true;
                self.error_message = None;
                self.generated_meme_url = None;

                let request = MemeRequest {
                    template: self.selected_template.clone(),
                    image_url: self.image_url.clone(),
                    top_text: self.top_text.clone(),
                    bottom_text: self.bottom_text.clone(),
                };

                let link = ctx.link().clone();
                spawn_local(async move {
                    match Request::post("/api/meme")
                        .json(&request)
                        .unwrap()
                        .send()
                        .await
                    {
                        Ok(response) => {
                            if response.ok() {
                                if let Ok(blob) = response.blob().await {
                                    let url = web_sys::Url::create_object_url_with_blob(&blob)
                                        .unwrap_or_default();
                                    link.send_message(Msg::MemeGenerated(url));
                                } else {
                                    link.send_message(Msg::Error("Failed to process meme image".to_string()));
                                }
                            } else {
                                link.send_message(Msg::Error("Failed to generate meme".to_string()));
                            }
                        }
                        Err(e) => {
                            link.send_message(Msg::Error(format!("Network error: {}", e)));
                        }
                    }
                });
                true
            }
            Msg::TemplatesLoaded(templates) => {
                self.templates = templates;
                true
            }
            Msg::MemeGenerated(url) => {
                self.generated_meme_url = Some(url);
                self.loading = false;
                true
            }
            Msg::Error(msg) => {
                self.error_message = Some(msg);
                self.loading = false;
                true
            }
        }
    }

    fn view(&self, ctx: &Context<Self>) -> Html {
        html! {
            <div class="container">
                <header>
                    <h1>{"🎭 Meme Generator"}</h1>
                    <p>{"Create hilarious memes with ease!"}</p>
                </header>

                <main>
                    <div class="form-section">
                        <h2>{"Choose Template or Image URL"}</h2>
                        
                        <div class="input-group">
                            <label>{"Select Template:"}</label>
                            <select
                                onchange={ctx.link().callback(|e: Event| {
                                    let target = e.target_dyn_into::<web_sys::HtmlSelectElement>();
                                    if let Some(select) = target {
                                        Msg::SelectTemplate(select.value())
                                    } else {
                                        Msg::SelectTemplate(String::new())
                                    }
                                })}
                                value={self.selected_template.clone()}
                            >
                                <option value="">{"-- Select a template --"}</option>
                                {
                                    self.templates.iter().map(|template| {
                                        html! {
                                            <option value={template.clone()}>{template}</option>
                                        }
                                    }).collect::<Html>()
                                }
                            </select>
                        </div>

                        <div class="separator">{"OR"}</div>

                        <div class="input-group">
                            <label>{"Image URL:"}</label>
                            <input
                                type="text"
                                placeholder="https://example.com/image.jpg"
                                value={self.image_url.clone()}
                                oninput={ctx.link().callback(|e: InputEvent| {
                                    let target = e.target_dyn_into::<HtmlInputElement>();
                                    if let Some(input) = target {
                                        Msg::UpdateImageUrl(input.value())
                                    } else {
                                        Msg::UpdateImageUrl(String::new())
                                    }
                                })}
                            />
                        </div>

                        <h2>{"Add Text"}</h2>
                        
                        <div class="input-group">
                            <label>{"Top Text:"}</label>
                            <input
                                type="text"
                                placeholder="Enter top text"
                                value={self.top_text.clone()}
                                oninput={ctx.link().callback(|e: InputEvent| {
                                    let target = e.target_dyn_into::<HtmlInputElement>();
                                    if let Some(input) = target {
                                        Msg::UpdateTopText(input.value())
                                    } else {
                                        Msg::UpdateTopText(String::new())
                                    }
                                })}
                            />
                        </div>

                        <div class="input-group">
                            <label>{"Bottom Text:"}</label>
                            <input
                                type="text"
                                placeholder="Enter bottom text"
                                value={self.bottom_text.clone()}
                                oninput={ctx.link().callback(|e: InputEvent| {
                                    let target = e.target_dyn_into::<HtmlInputElement>();
                                    if let Some(input) = target {
                                        Msg::UpdateBottomText(input.value())
                                    } else {
                                        Msg::UpdateBottomText(String::new())
                                    }
                                })}
                            />
                        </div>

                        <button
                            class="generate-btn"
                            onclick={ctx.link().callback(|_| Msg::GenerateMeme)}
                            disabled={self.loading || (self.selected_template.is_empty() && self.image_url.is_empty())}
                        >
                            {if self.loading { "Generating..." } else { "Generate Meme" }}
                        </button>
                    </div>

                    {
                        if let Some(error) = &self.error_message {
                            html! {
                                <div class="error-message">
                                    <p>{format!("❌ {}", error)}</p>
                                </div>
                            }
                        } else {
                            html! {}
                        }
                    }

                    {
                        if let Some(url) = &self.generated_meme_url {
                            html! {
                                <div class="result-section">
                                    <h2>{"Your Meme:"}</h2>
                                    <img src={url.clone()} alt="Generated Meme" />
                                    <a href={url.clone()} download="meme.png" class="download-btn">
                                        {"Download Meme"}
                                    </a>
                                </div>
                            }
                        } else {
                            html! {}
                        }
                    }
                </main>

                <footer>
                    <p>{"Built with ❤️ using Rust (Yew) and Go"}</p>
                </footer>
            </div>
        }
    }
}

fn main() {
    yew::Renderer::<App>::new().render();
}
