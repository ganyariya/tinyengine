#version 410 core

// 入力頂点属性
layout (location = 0) in vec3 aPos;      // 頂点位置
layout (location = 1) in vec3 aColor;    // 頂点色

// 出力変数
out vec3 vertexColor;

void main() {
    // 頂点位置をそのまま使用（変換なし）
    gl_Position = vec4(aPos, 1.0);
    
    // 頂点色をフラグメントシェーダーに渡す
    vertexColor = aColor;
}