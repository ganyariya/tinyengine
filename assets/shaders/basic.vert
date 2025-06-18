#version 410 core

// 入力頂点属性
layout (location = 0) in vec3 aPos;      // 頂点位置
layout (location = 1) in vec3 aColor;    // 頂点色

// 出力変数（フラグメントシェーダーに渡される）
out vec3 vertexColor;

// ユニフォーム変数
uniform mat4 model;         // モデル変換行列
uniform mat4 view;          // ビュー変換行列  
uniform mat4 projection;    // プロジェクション変換行列

void main() {
    // MVP変換を適用して最終的な頂点位置を計算
    gl_Position = projection * view * model * vec4(aPos, 1.0);
    
    // 頂点色をフラグメントシェーダーに渡す
    vertexColor = aColor;
}