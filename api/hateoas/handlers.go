package hateoas

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/juju/errors"
	"github.com/loopfz/gadgeto/tonic"
)

const hateoasBasePathKey = "HateoasBasePath"

// HandlerIndex generates a simple hateoas index
func HandlerIndex(links ...ResourceLink) gin.HandlerFunc {
	return tonic.Handler(func(c *gin.Context) (Resource, error) {
		var l []ResourceLink
		for _, link := range links {
			l = append(l, ResourceLink{
				Href: BaseURL(c) + link.Href,
				Rel:  link.Rel,
			})
		}
		return Resource{
			Links: l,
		}, nil
	}, http.StatusOK)
}

// HandlerFindByPage returns a filtered and paginated resource list
func HandlerFindByPage(repository PageableRepository) gin.HandlerFunc {
	return tonic.Handler(func(c *gin.Context) (*PagedResources, error) {
		pageable := Pageable{}
		criteria := map[string]interface{}{}
		if err := c.ShouldBindQuery(&pageable); err != nil {
			return nil, err
		}

		for k, v := range c.Request.URL.Query() {
			criteria[k] = v[0]
		}
		delete(criteria, "page")
		delete(criteria, "sort")
		delete(criteria, "size")
		delete(criteria, "indexedBy")
		results, err := repository.FindPageBy(pageable, criteria)
		if err != nil {
			return nil, err
		}

		resources := results.ToResources(BaseURL(c))
		return &resources, nil
	}, http.StatusPartialContent)
}

// HandlerFindBy returns all resources matching path params
func HandlerFindBy(repository ListableRepository) gin.HandlerFunc {
	return tonic.Handler(func(c *gin.Context) (interface{}, error) {
		return repository.FindBy(parsePathParams(c))
	}, http.StatusOK)
}

// HandlerFindOneBy returns the first resource matching path params
func HandlerFindOneBy(repository ListableRepository) gin.HandlerFunc {
	return tonic.Handler(func(c *gin.Context) (interface{}, error) {
		result, err := FindByPath(c, repository)
		if err != nil {
			return nil, err
		}
		return result, nil
	}, http.StatusOK)
}

// HandlerRemoveOneBy removes a given resource
func HandlerRemoveOneBy(repository SavableRepository) gin.HandlerFunc {
	return tonic.Handler(func(c *gin.Context) error {
		result, err := FindByPath(c, repository)
		if err != nil {
			return err
		}
		return repository.Remove(result)
	}, http.StatusNoContent)
}

// HandlerRemoveAll removes a whole collection
func HandlerRemoveAll(repository TruncatableRepository) gin.HandlerFunc {
	return tonic.Handler(func(c *gin.Context) error {
		return repository.Truncate()
	}, http.StatusNoContent)
}

// ErrorHook Convert repository errors in juju errors
func ErrorHook(tonicErrorHook tonic.ErrorHook) tonic.ErrorHook {
	return func(c *gin.Context, err error) (int, interface{}) {
		if errors.IsAlreadyExists(err) {
			return http.StatusConflict, nil
		}
		switch inner := err.(type) {
		case errorCreated:
			return http.StatusCreated, nil
		case errorGone:
			return http.StatusGone, nil
		case EntityDoesNotExistError:
			err = errors.NewNotFound(inner, inner.Error())
		case UnsupportedIndexError:
			err = errors.NewNotSupported(err, err.Error())
		}
		return tonicErrorHook(c, err)
	}
}

// AddToBasePath add a subpath to the BasePath stored in the gin context, in order to build hateoas links
func AddToBasePath(basePath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		path, ok := c.Get(hateoasBasePathKey)
		if !ok {
			c.Set(hateoasBasePathKey, basePath)
			c.Next()
			return
		}
		c.Set(hateoasBasePathKey, path.(string)+basePath)
	}
}

// FindByPath find one entity in the given repository, using paths parameters as matching criteria
func FindByPath(c *gin.Context, repository ListableRepository) (Entity, error) {
	params := parsePathParams(c)
	if repo, ok := repository.(SoftDeletableRepository); ok {
		result, err := repo.FindOneByUnscoped(params)
		if err != nil {
			return nil, err
		}

		if result.GetDeletedAt() != nil {
			return result, ErrorGone
		}
		return result, err
	}
	return repository.FindOneBy(params)
}

func parsePathParams(c *gin.Context) map[string]interface{} {
	criteria := map[string]interface{}{}
	for _, p := range c.Params {
		criteria[p.Key] = p.Value
	}
	return criteria
}
