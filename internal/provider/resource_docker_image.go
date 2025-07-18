package provider

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	dockerImageCreateDefaultTimeout = 20 * time.Minute
	dockerImageUpdateDefaultTimeout = 20 * time.Minute
	dockerImageDeleteDefaultTimeout = 20 * time.Minute
)

func resourceDockerImage() *schema.Resource {
	return &schema.Resource{
		Description: "Manages the lifecycle of a docker image in your docker host. It can be used to build a new docker image or to pull an existing one from a registry.\n This resource will *not* pull new layers of the image automatically unless used in conjunction with [docker_registry_image](../data-sources/registry_image.md) data source to update the `pull_triggers` field.",

		CreateContext: resourceDockerImageCreate,
		ReadContext:   resourceDockerImageRead,
		UpdateContext: resourceDockerImageUpdate,
		DeleteContext: resourceDockerImageDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(dockerImageCreateDefaultTimeout),
			Update: schema.DefaultTimeout(dockerImageUpdateDefaultTimeout),
			Delete: schema.DefaultTimeout(dockerImageDeleteDefaultTimeout),
		},

		Schema: map[string]*schema.Schema{
			"docker_client": dockerSchema,
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique identifier for this resource. This is not the image ID, but the ID of the resource in the Terraform state. This is used to identify the resource in the Terraform state. To reference the correct image ID, use the `image_id` attribute.",
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the Docker image, including any tags or SHA256 repo digests.",
				Required:    true,
				ForceNew:    true,
			},

			"image_id": {
				Type:        schema.TypeString,
				Description: "The ID of the image (as seen when executing `docker inspect` on the image). Can be used to reference the image via its ID in other resources.",
				Computed:    true,
			},

			"repo_digest": {
				Type:        schema.TypeString,
				Description: "The image sha256 digest in the form of `repo[:tag]@sha256:<hash>`. This may not be populated when building an image, because it is read from the local Docker client and so may be available only when the image was either pulled from the repo or pushed to the repo (perhaps using `docker_registry_image`) in a previous run.",
				Computed:    true,
			},

			"keep_locally": {
				Type:        schema.TypeBool,
				Description: "If true, then the Docker image won't be deleted on destroy operation. If this is false, it will delete the image from the docker local storage on destroy operation.",
				Optional:    true,
			},

			"pull_triggers": {
				Type:        schema.TypeSet,
				Description: "List of values which cause an image pull when changed. This is used to store the image digest from the registry when using the [docker_registry_image](../data-sources/registry_image.md).",
				Optional:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
			},

			"force_remove": {
				Type:        schema.TypeBool,
				Description: "If true, then the image is removed forcibly when the resource is destroyed.",
				Optional:    true,
			},

			"build": {
				Type:          schema.TypeSet,
				Description:   "Configuration to build an image. Requires the `Use containerd for pulling and storing images` option to be disabled in the Docker Host(https://github.com/kreuzwerker/terraform-provider-docker/issues/534). Please see [docker build command reference](https://docs.docker.com/engine/reference/commandline/build/#options) too.",
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"pull_triggers"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dockerfile": {
							Type:        schema.TypeString,
							Description: "Name of the Dockerfile. Defaults to `Dockerfile`.",
							Optional:    true,
							Default:     "Dockerfile",
							ForceNew:    true,
						},
						"tag": {
							Type:        schema.TypeList,
							Description: "Name and optionally a tag in the 'name:tag' format",
							Optional:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"remove": {
							Type:        schema.TypeBool,
							Description: "Remove intermediate containers after a successful build. Defaults to `true`.",
							Default:     true,
							Optional:    true,
						},
						"secrets": {
							Type:        schema.TypeList,
							Description: "Set build-time secrets. Only available when you use a buildx builder.",
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Description: "ID of the secret. By default, secrets are mounted to /run/secrets/<id>",
										Optional:    false,
										Required:    true,
										ForceNew:    true,
									},
									"src": {
										Type:        schema.TypeString,
										Description: "File source of the secret. Takes precedence over `env`",
										Optional:    true,
										Required:    false,
										ForceNew:    true,
									},
									"env": {
										Type:        schema.TypeString,
										Description: "Environment variable source of the secret",
										Optional:    true,
										Required:    false,
										ForceNew:    true,
									},
								},
							},
						},
						"label": {
							Type:        schema.TypeMap,
							Description: "Set metadata for an image",
							Optional:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"suppress_output": {
							Type:        schema.TypeBool,
							Description: "Suppress the build output and print image ID on success",
							Optional:    true,
							ForceNew:    true,
						},
						"remote_context": {
							Type:        schema.TypeString,
							Description: "A Git repository URI or HTTP/HTTPS context URI. Will be ignored if `builder` is set.",
							Optional:    true,
							ForceNew:    true,
						},
						"no_cache": {
							Type:        schema.TypeBool,
							Description: "Do not use the cache when building the image",
							Optional:    true,
							ForceNew:    true,
						},
						"force_remove": {
							Type:        schema.TypeBool,
							Description: "Always remove intermediate containers",
							Optional:    true,
							ForceNew:    true,
						},
						"pull_parent": {
							Type:        schema.TypeBool,
							Description: "Attempt to pull the image even if an older image exists locally",
							Optional:    true,
							ForceNew:    true,
						},
						"isolation": {
							Type:        schema.TypeString,
							Description: "Isolation represents the isolation technology of a container. The supported values are ",
							Optional:    true,
							ForceNew:    true,
						},
						"cpu_set_cpus": {
							Type:        schema.TypeString,
							Description: "CPUs in which to allow execution (e.g., `0-3`, `0`, `1`)",
							Optional:    true,
							ForceNew:    true,
						},
						"cpu_set_mems": {
							Type:        schema.TypeString,
							Description: "MEMs in which to allow execution (`0-3`, `0`, `1`)",
							Optional:    true,
							ForceNew:    true,
						},
						"cpu_shares": {
							Type:        schema.TypeInt,
							Description: "CPU shares (relative weight)",
							Optional:    true,
							ForceNew:    true,
						},
						"cpu_quota": {
							Type:        schema.TypeInt,
							Description: "Microseconds of CPU time that the container can get in a CPU period",
							Optional:    true,
							ForceNew:    true,
						},
						"cpu_period": {
							Type:        schema.TypeInt,
							Description: "The length of a CPU period in microseconds",
							Optional:    true,
							ForceNew:    true,
						},
						"memory": {
							Type:        schema.TypeInt,
							Description: "Set memory limit for build",
							Optional:    true,
							ForceNew:    true,
						},
						"memory_swap": {
							Type:        schema.TypeInt,
							Description: "Total memory (memory + swap), -1 to enable unlimited swap",
							Optional:    true,
							ForceNew:    true,
						},
						"cgroup_parent": {
							Type:        schema.TypeString,
							Description: "Optional parent cgroup for the container",
							Optional:    true,
							ForceNew:    true,
						},
						"network_mode": {
							Type:        schema.TypeString,
							Description: "Set the networking mode for the RUN instructions during build",
							Optional:    true,
							ForceNew:    true,
						},
						"shm_size": {
							Type:        schema.TypeInt,
							Description: "Size of /dev/shm in bytes. The size must be greater than 0",
							Optional:    true,
							ForceNew:    true,
						},
						"ulimit": {
							Type:        schema.TypeList,
							Description: "Configuration for ulimits",
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Description: "type of ulimit, e.g. `nofile`",
										Required:    true,
										ForceNew:    true,
									},
									"hard": {
										Type:        schema.TypeInt,
										Description: "soft limit",
										Required:    true,
										ForceNew:    true,
									},
									"soft": {
										Type:        schema.TypeInt,
										Description: "hard limit",
										Required:    true,
										ForceNew:    true,
									},
								},
							},
						},
						"build_args": {
							Type:        schema.TypeMap,
							Description: "Pairs for build-time variables in the form of `ENDPOINT : \"https://example.com\"`",
							Optional:    true,
							ForceNew:    true,
							Elem: &schema.Schema{
								Type:        schema.TypeString,
								Description: "The argument",
							},
						},
						"auth_config": {
							Type:        schema.TypeList,
							Description: "The configuration for the authentication",
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"host_name": {
										Type:        schema.TypeString,
										Description: "hostname of the registry",
										Required:    true,
									},
									"user_name": {
										Type:        schema.TypeString,
										Description: "the registry user name",
										Optional:    true,
									},
									"password": {
										Type:        schema.TypeString,
										Description: "the registry password",
										Optional:    true,
									},
									"auth": {
										Type:        schema.TypeString,
										Description: "the auth token",
										Optional:    true,
									},
									"email": {
										Type:        schema.TypeString,
										Description: "the user emal",
										Optional:    true,
									},
									"server_address": {
										Type:        schema.TypeString,
										Description: "the server address",
										Optional:    true,
									},
									"identity_token": {
										Type:        schema.TypeString,
										Description: "the identity token",
										Optional:    true,
									},
									"registry_token": {
										Type:        schema.TypeString,
										Description: "the registry token",
										Optional:    true,
									},
								},
							},
						},
						"context": {
							Type:        schema.TypeString,
							Description: "Value to specify the build context. Currently, only a `PATH` context is supported. You can use the helper function '${path.cwd}/context-dir'. This always refers to the local working directory, even when building images on remote hosts. Please see https://docs.docker.com/build/building/context/ for more information about build contexts.",
							Required:    true,
							ForceNew:    true,
						},
						"labels": {
							Type:        schema.TypeMap,
							Description: "User-defined key/value metadata",
							Optional:    true,
							ForceNew:    true,
							Elem: &schema.Schema{
								Type:        schema.TypeString,
								Description: "The key/value pair",
							},
						},
						"squash": {
							Type:        schema.TypeBool,
							Description: "If true the new layers are squashed into a new image with a single new layer",
							Optional:    true,
							ForceNew:    true,
						},
						"cache_from": {
							Type:        schema.TypeList,
							Description: "Images to consider as cache sources",
							Optional:    true,
							ForceNew:    true,
							Elem: &schema.Schema{
								Type:        schema.TypeString,
								Description: "The image",
							},
						},
						"security_opt": {
							Type:        schema.TypeList,
							Description: "The security options",
							Optional:    true,
							ForceNew:    true,
							Elem: &schema.Schema{
								Type:        schema.TypeString,
								Description: "The option",
							},
						},
						"extra_hosts": {
							Type:        schema.TypeList,
							Description: "A list of hostnames/IP mappings to add to the container’s /etc/hosts file. Specified in the form [\"hostname:IP\"]",
							Optional:    true,
							ForceNew:    true,
							Elem: &schema.Schema{
								Type:        schema.TypeString,
								Description: "",
							},
						},
						"target": {
							Type:        schema.TypeString,
							Description: "Set the target build stage to build",
							Optional:    true,
							ForceNew:    true,
						},
						"session_id": {
							Type:        schema.TypeString,
							Description: "Set an ID for the build session",
							Optional:    true,
							ForceNew:    true,
						},
						"platform": {
							Type:        schema.TypeString,
							Description: "Set the target platform for the build. Defaults to `GOOS/GOARCH`. For more information see the [docker documentation](https://github.com/docker/buildx/blob/master/docs/reference/buildx.md#-set-the-target-platforms-for-the-build---platform)",
							Optional:    true,
							ForceNew:    true,
						},
						"version": {
							Type:        schema.TypeString,
							Description: "Version of the underlying builder to use",
							Optional:    true,
							ForceNew:    true,
						},
						"build_id": {
							Type:        schema.TypeString,
							Description: "BuildID is an optional identifier that can be passed together with the build request. The same identifier can be used to gracefully cancel the build with the cancel request.",
							Optional:    true,
							ForceNew:    true,
						},
						"builder": {
							Type:        schema.TypeString,
							Description: "Set the name of the buildx builder to use. If not set, the legacy builder is used.",
							Optional:    true,
							ForceNew:    true,
						},
						"build_log_file": {
							Type:        schema.TypeString,
							Description: "Path to a file where the buildx log are written to. Only available when `builder` is set. If not set, no logs are available. The path is taken as is, so make sure to use a path that is available.",
							Optional:    true,
							ForceNew:    true,
						},
					},
				},
			},
			"triggers": {
				Description: "A map of arbitrary strings that, when changed, will force the `docker_image` resource to be replaced. This can be used to rebuild an image when contents of source code folders change",
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
			},
			"platform": {
				Type:        schema.TypeString,
				Description: "The platform to use when pulling the image. Defaults to the platform of the current machine.",
				Optional:    true,
				Default:     "",
				ForceNew:    true,
			},
		},
	}
}
