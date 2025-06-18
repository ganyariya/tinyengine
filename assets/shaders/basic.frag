#version 410 core

// 入力変数（頂点シェーダーから受け取る）
in vec3 vertexColor;

// 出力変数（最終的なピクセル色）
out vec4 FragColor;

// ユニフォーム変数（オプション）
uniform float alpha;        // アルファ値（透明度）
uniform float time;         // 時間（アニメーション用）

void main() {
    // 基本的な色の出力
    FragColor = vec4(vertexColor, alpha);
    
    // 時間ベースの色変化効果（オプション）
    // フラグメント色に時間ベースの効果を追加
    float pulseFactor = 0.5 + 0.5 * sin(time * 2.0);
    FragColor.rgb *= pulseFactor;
}