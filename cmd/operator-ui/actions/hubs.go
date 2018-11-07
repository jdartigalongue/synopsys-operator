package actions

import (
	"fmt"
	"strings"

	"github.com/blackducksoftware/perceptor-protoform/pkg/api/hub/v1"
	"github.com/blackducksoftware/perceptor-protoform/pkg/hub"
	hubclientset "github.com/blackducksoftware/perceptor-protoform/pkg/hub/client/clientset/versioned"
	"github.com/blackducksoftware/perceptor-protoform/pkg/util"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/validate"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// This file is generated by Buffalo. It offers a basic structure for
// adding, editing and deleting a page. If your model is more
// complex or you need more than the basic implementation you need to
// edit this file.

// Following naming logic is implemented in Buffalo:
// Model: Singular (Hub)
// DB Table: Plural (Hubs)
// Resource: Plural (Hubs)
// Path: Plural (/hubs)
// View Template Folder: Plural (/templates/hubs/)

// HubsResource is the resource for the Blackduck model
type HubsResource struct {
	buffalo.Resource
	kubeClient *kubernetes.Clientset
	hubClient  *hubclientset.Clientset
}

// NewHubResource will instantiate the Black Duck Resource
func NewHubResource(kubeConfig *rest.Config) (*HubsResource, error) {
	kubeClient, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to create kube client due to %+v", err)
	}
	hubClient, err := hubclientset.NewForConfig(kubeConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to create hub client due to %+v", err)
	}
	return &HubsResource{kubeClient: kubeClient, hubClient: hubClient}, nil
}

// List gets all Hubs. This function is mapped to the path
// GET /hubs
func (v HubsResource) List(c buffalo.Context) error {
	blackducks, _ := util.ListHubs(v.hubClient, "")
	// Make blackducks available inside the html template
	c.Set("hubs", blackducks.Items)
	return c.Render(200, r.HTML("hubs/index.html", "old_application.html"))
}

// Show gets the data for one Hub. This function is mapped to
// the path GET /hubs/{hub_id}
func (v HubsResource) Show(c buffalo.Context) error {
	blackduck, _ := util.GetHub(v.hubClient, c.Param("hub_id"), c.Param("hub_id"))
	// Make blackduck available inside the html template
	c.Set("hub", blackduck)
	return c.Render(200, r.HTML("hubs/show.html", "old_application.html"))
}

// New renders the form for creating a new Hub.
// This function is mapped to the path GET /hubs/new
func (v HubsResource) New(c buffalo.Context) error {
	blackduck := &v1.Hub{}
	err := v.common(c, blackduck)
	if err != nil {
		return err
	}
	// Make blackduck available inside the html template
	c.Set("hub", blackduck)

	return c.Render(200, r.HTML("hubs/new.html", "old_application.html"))
}

func (v HubsResource) common(c buffalo.Context, blackduck *v1.Hub) error {
	var storageList map[string]string
	storageList = make(map[string]string)
	storageClasses, err := util.ListStorageClass(v.kubeClient)
	if err != nil {
		c.Error(404, fmt.Errorf("\"message\": \"Failed to List the storage class due to %+v\"", err))
	}
	for _, storageClass := range storageClasses.Items {
		storageList[fmt.Sprintf("%s (%s)", storageClass.GetName(), storageClass.Provisioner)] = storageClass.GetName()
	}
	storageList[fmt.Sprintf("%s (%s)", "None", "Disable dynamic provisioner")] = "none"
	blackduck.View.StorageClasses = storageList

	keys, _ := util.ListHubPV(v.hubClient, "")
	blackduck.View.Clones = keys

	blackducks, _ := util.ListHubs(v.hubClient, "")
	certificateNames := []string{"default", "manual"}
	for _, hub := range blackducks.Items {
		if strings.EqualFold(hub.Spec.CertificateName, "manual") {
			certificateNames = append(certificateNames, hub.Spec.Namespace)
		}
	}
	blackduck.View.CertificateNames = certificateNames
	env, images := hub.GetHubKnobs()
	// env := map[string]string{}
	// images := []string{}
	environs := []string{}
	for key, value := range env {
		if !strings.EqualFold(value, "") {
			environs = append(environs, fmt.Sprintf("%s:%s", key, value))
		}
	}

	// environs := []string{"IPV4_ONLY:0", "HUB_PROXY_NON_PROXY_HOSTS:solr"}
	blackduck.View.Environs = environs

	blackduck.View.ContainerTags = images
	return nil
}

func (v HubsResource) redirect(c buffalo.Context, blackduck *v1.Hub, err error) error {
	if err != nil {
		// Make blackduck available inside the html template
		err := v.common(c, blackduck)
		if err != nil {
			log.Error(err)
			return err
		}
		log.Infof("edit hub in create: %+v", blackduck)
		c.Set("hub", blackduck)
		log.Info("Before")
		validateErrs := validate.NewErrors()
		log.Info("After")
		// validateErrs.Add("error", err.Error())
		log.Infof("validateErrs: %+v", validateErrs)
		// validateErrs.Errors = map[string][]string{"error": []string{errors.WithStack(err).Error()}}
		c.Set("errors", validateErrs)
		return c.Render(422, r.HTML("hubs/new.html", "old_application.html"))
	}
	return nil
}

// Create adds a Blackduck to the DB. This function is mapped to the
// path POST /hubs
func (v HubsResource) Create(c buffalo.Context) error {
	// Allocate an empty Blackduck
	hub := &v1.Hub{}

	// Bind blackduck to the html form elements
	if err := c.Bind(hub); err != nil {
		log.Errorf("unable to bind blackduck %+v because %+v", c, err)
		return errors.WithStack(err)
	}

	log.Infof("create hub: %+v", hub)

	ns, err := util.CreateNamespace(v.kubeClient, hub.Spec.Namespace)
	if err != nil {
		v.redirect(c, hub, err)
	}
	log.Infof("created namespace for %s is %+v", hub.Spec.Namespace, ns)

	_, err = util.CreateHub(v.hubClient, hub.Spec.Namespace, &v1.Hub{ObjectMeta: metav1.ObjectMeta{Name: hub.Spec.Namespace}, Spec: hub.Spec})

	if err != nil {
		v.redirect(c, hub, err)
	}
	// If there are no errors set a success message
	c.Flash().Add("success", "Blackduck was created successfully")

	blackducks, _ := util.ListHubs(v.hubClient, "")
	c.Set("hubs", blackducks.Items)
	// and redirect to the blackducks index page
	return c.Redirect(302, "/hubs/%s", hub.Spec.Namespace)
}

// Edit renders a edit form for a Hub. This function is
// mapped to the path GET /hubs/{hub_id}/edit
func (v HubsResource) Edit(c buffalo.Context) error {
	// Get the DB connection from the context
	// tx, ok := c.Value("tx").(*pop.Connection)
	// if !ok {
	// 	return errors.WithStack(errors.New("no transaction found"))
	// }

	// // Allocate an empty Blackduck
	// blackduck := &v1.Hub{}

	// if err := tx.Find(blackduck, c.Param("blackduck_id")); err != nil {
	// 	return c.Error(404, err)
	// }

	// return c.Render(200, r.Auto(c, blackduck))
	return c.Error(404, errors.New("resource not implemented"))

}

// Update changes a Blackduck in the DB. This function is mapped to
// the path PUT /hubs/{hub_id}
func (v HubsResource) Update(c buffalo.Context) error {
	// Get the DB connection from the context
	// tx, ok := c.Value("tx").(*pop.Connection)
	// if !ok {
	// 	return errors.WithStack(errors.New("no transaction found"))
	// }

	// // Allocate an empty Blackduck
	// blackduck := &v1.Hub{}

	// if err := tx.Find(blackduck, c.Param("blackduck_id")); err != nil {
	// 	return c.Error(404, err)
	// }

	// // Bind Blackduck to the html form elements
	// if err := c.Bind(blackduck); err != nil {
	// 	return errors.WithStack(err)
	// }

	// verrs, err := tx.ValidateAndUpdate(blackduck)
	// if err != nil {
	// 	return errors.WithStack(err)
	// }

	// if verrs.HasAny() {
	// 	// Make the errors available inside the html template
	// 	c.Set("errors", verrs)

	// 	// Render again the edit.html template that the user can
	// 	// correct the input.
	// 	return c.Render(422, r.Auto(c, blackduck))
	// }

	// // If there are no errors set a success message
	// c.Flash().Add("success", "Blackduck was updated successfully")

	// and redirect to the blackducks index page
	return c.Error(404, errors.New("resource not implemented"))
}

// Destroy deletes a Hub from the DB. This function is mapped
// to the path DELETE /hubs/{hub_id}
func (v HubsResource) Destroy(c buffalo.Context) error {

	log.Infof("delete hub request %v", c.Param("hub_id"))

	_, err := util.GetHub(v.hubClient, c.Param("hub_id"), c.Param("hub_id"))
	// To find the Blackduck the parameter blackduck_id is used.
	if err != nil {
		return c.Error(404, err)
	}

	// This is on the event loop.
	err = v.hubClient.SynopsysV1().Hubs(c.Param("hub_id")).Delete(c.Param("hub_id"), &metav1.DeleteOptions{})

	// To find the Blackduck the parameter blackduck_id is used.
	if err != nil {
		return c.Error(404, err)
	}

	// If there are no errors set a flash message
	c.Flash().Add("success", "Blackduck was deleted successfully")

	// blackducks, _ := util.ListHubs(v.hubClient, "")
	// c.Set("hubs", blackducks.Items)

	// Redirect to the blackducks index page
	return c.Redirect(302, "/hubs")
}
