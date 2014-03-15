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
out vec4 fpos_eye; // Fragment position in eye space.  For View vector.
out vec4 fnor_eye;
out vec4 fnor_world; // Used to point to environment maps.
out vec2 fuv;
void main(){
	fpos_eye = model_to_eye * vec4(vpos, 1.0);
	fnor_eye = model_to_eye * vec4(vnor, 0.0);
    fnor_world = eye_to_world * fnor_eye;
	fuv = vuv;
	gl_Position = eye_to_clip * fpos_eye;
}`},
	FSH_NOR_UV: ShaderSeed{Type: FRAGMENT_SHADER, Source: `
#version 330 core
// Takes normal from vertex.
in vec4 fpos_eye;
in vec4 fnor_eye;
in vec4 fnor_world; // Environment map is in world space.
in vec2 fuv;
uniform samplerCube environment_map;
uniform sampler2D albedo_map;
uniform sampler2D normal_map;
out vec3 color;

#define PI       3.1415926535897932384626433832795
#define TAU      6.2831853071795864769252867665590
#define SQRT_TAU 2.5066282746310005024157652848110

// Schlick approximation of Fresnel's reflectance.
// cspec is a specular color defined for normal incidence.
//     Typically 2..5 % for dielectrics, and 50..100 % for metals.
//     Dielectrics: r=g=b.  In metals, it varies.
//     Note that I don't see any reason to let alpha out; I treat it as a color
//     until I have another use for it.  Just say it's UV.
// l is the direction of the light.  It is pointing out of the surface.
// h is half way between the direction of the light and that of the view.
//     It makes sense because we are looking at the microfacets that reflect the
//     light into our eyes, and these microfacets have a normal h.
// Here, l and h must be in the same reference frame, it does not matter which
// one.
vec4 fresnel(vec4 cspec, vec3 l, vec3 h) {
	float cosangle = max(0, dot(l, h));
    return cspec + (1 - cspec ) * pow(1 - cosangle, 5);
}

// Normal distribution term of microfacets.
// h is a direction pointing out of the surface.
//     Currently expressed in the eye reference frame.  This may have to change
//     when we start caring about anisotropy and/or normal mapping.
float mfNormalDist(float sigma, vec3 v, vec3 h) {
    // This is a "normal distribution".  Not to be confused with the normal
	// distribution.  Yeah, same names for two things, amazing.  This whole
	// function computes a distribution of normal vectors.  To do this, it
	// uses a continuous probability distribution called "normal distribution".
    // It's also called "gaussian".
    float x = dot(h, v) - 1;
	float si = 1 / sigma;
	return SQRT_TAU * si * exp(-x*x*si*si*.5);
}

void main(){
	vec3 normal_world = normalize(fnor_world);
    // Surface albedo.
	vec3 albedo = vec3(.1, .1, .1);
    // Hardcoded sun light.
	vec3 l_dir = normalize(vec3(1, -5, 10)); // World coordinates for now.
    vec3 l_col = vec3(1, 1, .8);
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
	color = albedo * (
	    textureLod(environment_map, normal_world, 1000).rgb +
        l_col * max(dot(l_dir, normal_world), 0)
	);
	// Tone mapping.
	color = color / (1+color);
}`},
}
