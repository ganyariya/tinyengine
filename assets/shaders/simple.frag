#version 410 core

// 入力変数
in vec3 vertexColor;

// 出力変数
out vec4 FragColor;

void main() {
    // 頂点色をそのまま出力
    FragColor = vec4(vertexColor, 1.0);
}