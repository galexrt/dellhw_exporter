The dellhw_exporter can be configured to cache the results to prevent unnecessary executions of `omreport`, which can lead to high load.

## Caching

Why do we even need caching in an exporter?

Unfortunately Prometheus itself does not provide a 'real' High-Availability concept. The Prometheus authors recommend running identical Prometheus instances to achieve HA (see [Prometheus FAQ](https://prometheus.io/docs/introduction/faq/#can-prometheus-be-made-highly-available)), hence the exporters will be scraped multiple times within one scrape interval.

Besides the problem, that not all instances are retrieving the identical metric values, this can also produce high load if the actual collection of metrics is 'expensive'. While the first problem is not a real use case, because Prometheus does not claim to be consistent, the second problem is a real problem and valid use case for caching.

In particular the `dellhw_exporter`, since the underlying `omreport` calls produce high load. This is caused by many drivers collecting data from different components.

### Configuration

As you may have seen in [the Configuration doc page](docs/configuration.md) there are two caching related configuration parameters for enablement and how long the cache should be valid.

```console
--cache-enabled bool           Enable caching (default false)
--cache-duration int           Duration in seconds for the cache lifetime (default 20)
```

If you want to retrieve new metrics on each scrape, but want to prevent multiple collections because of multiple Prometheus instances, it is a good idea to set the cache-duration equal to your job's `scrape_interval`. If the `scrape_interval` is even less than the default value it can be useful to set a different `cache-duration`, maybe 2-3 times of the `scrape_interval`.


### Implementation details

An additional adapter channel is used to retrieve the collected metrics and put them into an array if caching is enabled. A mutex is used to prevent concurrent collections and concurrent write operations to the cache.

Since the metrics are pushed into a "local" channel instead of the channel passed by the Prometheus library directly, we need a second waitgroup. The first waitgroup ensures that all collectors have finished. The second waitgroup ensures that all metrics are written to the outgoing channel before the method returns. This is needed because the Prometheus library will close the channel once the method returns.

For further details please see the initial [Caching PR](https://github.com/galexrt/dellhw_exporter/pull/46).
