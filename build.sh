#!/bin/bash

echo "🚀 Building Snippet Manager..."

# Build Svelte frontend
echo "📦 Building frontend..."
cd web || exit
yarn build
cd ..

# Copy build to server directory for embed
echo "📂 Copying static files..."
rm -rf internal/server/dist
cp -r web/build internal/server/dist

# Build Go backend
echo "🔨 Building backend..."
go build -o sni ./cmd/sni

echo "✅ Build complete!"
echo ""
echo "Usage:"
echo "  ./sni --help           - Show help"
echo "  ./sni list             - List all snippets"  
echo "  ./sni new <name>       - Create new snippet"
echo "  ./sni server           - Start web UI server"
echo ""
echo "🌐 Web UI will be available at http://localhost:8080"
