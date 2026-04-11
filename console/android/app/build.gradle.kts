plugins {
    id("com.android.application")
    id("kotlin-android")
    // The Flutter Gradle Plugin must be applied after the Android and Kotlin Gradle plugins.
    id("dev.flutter.flutter-gradle-plugin")
}

android {
    namespace = "space.retrovibe.retrovibed"
    compileSdk = flutter.compileSdkVersion
    ndkVersion = flutter.ndkVersion

    compileOptions {
        sourceCompatibility = JavaVersion.VERSION_17
        targetCompatibility = JavaVersion.VERSION_17
    }

    kotlinOptions {
        jvmTarget = JavaVersion.VERSION_17.toString()
    }

    signingConfigs {
        create("release") {
            storeFile = System.getenv("RETROVIBED_ANDROID_KEY_STORE_PATH")?.let { file(it) }
            storePassword = System.getenv("RETROVIBED_ANDROID_STORE_PASSWORD")
            keyAlias = System.getenv("RETROVIBED_ANDROID_KEY_ALIAS")
            keyPassword = System.getenv("RETROVIBED_ANDROID_STORE_PASSWORD")
        }
    }

    defaultConfig {
        applicationId = "space.retrovibe.retrovibed"
        minSdk = flutter.minSdkVersion
        targetSdk = flutter.targetSdkVersion
        versionCode = flutter.versionCode
        versionName = flutter.versionName
    }

    buildTypes {
        release {
            signingConfig = signingConfigs.getByName("release")
        }
    }
}

flutter {
    source = "../.."
}
