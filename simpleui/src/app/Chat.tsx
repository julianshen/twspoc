"use client";
import React, { useState, useRef } from "react";

type Message = {
  role: "user" | "assistant";
  content: string;
};

const OLLAMA_URL = "http://localhost:11434/api/chat";
const MODEL = "gemma3-4b";

export default function Chat() {
  const [messages, setMessages] = useState<Message[]>([]);
  const [input, setInput] = useState("");
  const [loading, setLoading] = useState(false);
  const messagesEndRef = useRef<HTMLDivElement>(null);

  // Scroll to bottom on new message
  React.useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages]);

  async function sendMessage(e: React.FormEvent) {
    e.preventDefault();
    if (!input.trim()) return;

    const newMessages: Message[] = [
      ...messages,
      { role: "user", content: input } as Message,
    ];
    setMessages(newMessages);
    setInput("");
    setLoading(true);

    try {
      const res = await fetch(OLLAMA_URL, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          model: MODEL,
          messages: newMessages,
          stream: false, // set to true for streaming
        }),
      });

      if (!res.ok) {
        throw new Error("Ollama API error");
      }

      const data = await res.json();
      const aiMessage = data.message?.content || data.content || "[No response]";
      setMessages((msgs) => [
        ...msgs,
        { role: "assistant", content: aiMessage } as Message,
      ]);
    } catch (err: unknown) {
      let errorMsg = "Unknown error";
      if (err instanceof Error) {
        errorMsg = err.message;
      }
      setMessages((msgs) => [
        ...msgs,
        { role: "assistant", content: "Error: " + errorMsg } as Message,
      ]);
    } finally {
      setLoading(false);
    }
  }

  return (
    <div style={{
      maxWidth: 600,
      margin: "0 auto",
      padding: 24,
      border: "1px solid #ddd",
      borderRadius: 8,
      background: "#fafafa"
    }}>
      <h2>AI Chat (Ollama: gemma3-4b)</h2>
      <div style={{
        minHeight: 300,
        maxHeight: 400,
        overflowY: "auto",
        marginBottom: 16,
        padding: 8,
        background: "#fff",
        border: "1px solid #eee",
        borderRadius: 4
      }}>
        {messages.map((msg, i) => (
          <div key={i} style={{
            textAlign: msg.role === "user" ? "right" : "left",
            margin: "8px 0"
          }}>
            <span style={{
              display: "inline-block",
              padding: "8px 12px",
              borderRadius: 16,
              background: msg.role === "user" ? "#d1eaff" : "#e6e6e6",
              color: "#222",
              maxWidth: "80%",
              wordBreak: "break-word"
            }}>
              {msg.content}
            </span>
          </div>
        ))}
        <div ref={messagesEndRef} />
      </div>
      <form onSubmit={sendMessage} style={{ display: "flex", gap: 8 }}>
        <input
          type="text"
          value={input}
          onChange={(e) => setInput(e.target.value)}
          placeholder="Type your message..."
          style={{
            flex: 1,
            padding: 10,
            borderRadius: 8,
            border: "1px solid #ccc"
          }}
          disabled={loading}
        />
        <button
          type="submit"
          disabled={loading || !input.trim()}
          style={{
            padding: "0 20px",
            borderRadius: 8,
            border: "none",
            background: "#0070f3",
            color: "#fff",
            fontWeight: "bold",
            cursor: loading ? "not-allowed" : "pointer"
          }}
        >
          {loading ? "..." : "Send"}
        </button>
      </form>
      <div style={{ fontSize: 12, color: "#888", marginTop: 8 }}>
        Ollama must be running locally with the <b>gemma3-4b</b> model pulled.
      </div>
    </div>
  );
}
