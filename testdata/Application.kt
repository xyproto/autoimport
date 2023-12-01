package com.ostekake.trust.feedback

import org.springframework.boot.SpringApplication
import org.springframework.boot.autoconfigure.EnableAutoConfiguration
import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.boot.autoconfigure.domain.EntityScan
import org.springframework.boot.autoconfigure.freemarker.FreeMarkerAutoConfiguration
import org.springframework.boot.autoconfigure.web.servlet.error.ErrorMvcAutoConfiguration
import java.util.*
import javax.annotation.PostConstruct

@SpringBootApplication(
    scanBasePackages = ["com.ostekake.trust.feedback"],
    scanBasePackageClasses = [PackageMarker::class]
)
@EnableAutoConfiguration(
    exclude = [
        FreeMarkerAutoConfiguration::class,
        ErrorMvcAutoConfiguration::class]
)
@EntityScan(basePackages = ["com.ostekake.trust.feedback"])
class OstekakeApplication {
    @PostConstruct
    fun setTimeZoneAtStartup() {
        TimeZone.setDefault(TimeZone.getTimeZone("UTC"))
    }

    companion object {
        @JvmStatic
        fun main(args: Array<String>) {
            SpringApplication.run(OstekakeApplication::class.java, *args)
        }
    }
}
