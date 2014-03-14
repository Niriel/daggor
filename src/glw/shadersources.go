package glw

var SHADER_SOURCES = map[ShaderRef]ShaderSeed{
	VSH_NOR_UV_INSTANCED: ShaderSeed{Type: VERTEX_SHADER, Source: `
#version 330 core
// Transmits the normal to the fragment shader.
// The model matrix is not in a uniform.  This shader is
// to be used with instanced rendering.

layout(location = 0) in vec3 vpos;
layout(location = 1) in vec3 vnor;
layout(location = 2) in vec2 vuv;
layout(location = 3) in mat4 model_to_eye; // Instanced attribute.
// Note that the model_to_eye matrix occupies 4 attribute positions.
// The next layout location would be 7.

layout(std140) uniform GlobalMatrices
{
    mat4 eye_to_clip;
    mat4 eye_to_world;
};
out vec3 fnor_eye;
out vec3 fnor_world; // Used to point to environment maps.
out vec2 fuv;
void main(){
	vec4 eyePos = model_to_eye * vec4(vpos, 1.0);
	gl_Position = eye_to_clip * eyePos;
	fnor_eye = (model_to_eye * vec4(vnor, 0.0)).xyz;
    fnor_world = (eye_to_world * vec4(fnor_eye, 0.0)).xyz;
	fuv = vuv;
}`},
	FSH_NOR_UV: ShaderSeed{Type: FRAGMENT_SHADER, Source: `
#version 330 core
// Takes normal from vertex.
in vec3 fnor_eye;
in vec3 fnor_world; // Environment map is in world space.
in vec2 fuv;
uniform samplerCube environment_map;
uniform sampler2D albedo_map;
uniform sampler2D normal_map;
out vec3 color;
void main(){
	vec3 normal_world = normalize(fnor_world);
	// The 1000 here is an insanely high Level Of Detail which will fall down
	// to the blurriest version of the texture there is.  This simulates the
	// lambertian reflection model of a surface: the ray from your eye is
	// reflected on the entire half-space that the surface sees, therefore you
	// see an average of that.
	// In other words: this is our ambiant lighting.  Instead of taking a single
	// color, it takes its value from the environment map.  This actually makes
	// perfect sense.  At sunset, a white building will then seem orangeish on
	// one side.  And a piece of paper held horizontally above a grass field
	// will be greenish on the lower side, and white (sun+sky) on the upper
	// side.
	color = textureLod(environment_map, normal_world, 1000).rgb;

	// Now, we must include an albedo: not everything reflects 100% of what it
	// receives.  The albedo can be a single number or come from an albedo
	// texture, which is akin to the diffuse texture in other games.
    // Here is a table of albedos:
	// Fresh asphalt 	0.04
	// Worn asphalt 	0.12
	// Conifer forest (Summer) 	0.08, 0.09 to 0.15
	// Deciduous trees 	0.15 to 0.18
	// Bare soil 	0.17
	// Green grass 	0.25
	// Desert sand 	0.40
	// New concrete 	0.55
	// Ocean ice 	0.5--0.7
	// Fresh snow 	0.80--0.90
    // Note that these albedos code the luminosity only, not the hue.
}`},
}
