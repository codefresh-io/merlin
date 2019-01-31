package cache

import (
	"fmt"
	"io/ioutil"
	"time"

	goCache "github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

const (
	defaultExpiration = 30 * time.Minute
	defaultCleanup    = 8 * time.Hour
)

type (
	Cache interface {
		Put(key string, value interface{})
		Get(key string, target interface{}) error
		Persist()
	}

	cache struct {
		client *goCache.Cache
		path   string
		log    *logrus.Entry
	}

	Options struct {
		Path    string
		Logger  *logrus.Entry
		NoCache bool
	}
)

func New(opt *Options) Cache {
	c := &cache{
		path: opt.Path,
		log:  opt.Logger,
	}
	if opt.NoCache {
		c.client = goCache.New(defaultExpiration, defaultCleanup)
		c.log.Debug("Skipping usage of cache")
		return c
	}
	res, err := ioutil.ReadFile(opt.Path)
	if err != nil {
		c.log.Debugf("Failed to load cache: %s", err.Error())
		c.log.Debug("Creating fresh cache")
		c.client = goCache.New(defaultExpiration, defaultCleanup)
	}

	items := map[string]goCache.Item{}
	err = yaml.Unmarshal(res, items)
	if err != nil {
		c.log.Debugf("Failed to unmarshal cache into objects: %s", err.Error())
		c.log.Debug("Creating fresh cache")
		c.client = goCache.New(defaultExpiration, defaultCleanup)
	}
	c.log.Debugf("Loaded from cache")
	c.client = goCache.NewFrom(defaultExpiration, defaultCleanup, items)
	return c
}

func (c *cache) Put(key string, value interface{}) {
	c.client.Set(key, value, defaultExpiration)
}
func (c *cache) Get(key string, target interface{}) error {
	foo, found := c.client.Get(key)
	if found {
		return c.convert(foo, target)
	}
	return fmt.Errorf("Not found")
}
func (c *cache) Persist() {
	c.log.Debugf("Persisting cache to %s", c.path)
	items := c.client.Items()
	res, err := yaml.Marshal(items)
	if err != nil {
		c.log.Debugf("Error occured: %s", err.Error())
		return
	}
	err = ioutil.WriteFile(c.path, res, 0644)
	if err != nil {
		c.log.Debugf("Error occured: %s", err.Error())
		return
	}
	c.log.Debug("Saved")

}

func (c *cache) convert(key interface{}, target interface{}) error {
	bytes, err := yaml.Marshal(key)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(bytes, target)
}
