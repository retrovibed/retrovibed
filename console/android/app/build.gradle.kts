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
            // Explicit here (matching AGP 9's actual default for release) so
            // proguard-rules.pro is guaranteed to apply. Needed to patch a
            // gap in mobile_scanner's consumer rules that otherwise lets R8
            // strip ML Kit classes it needs at runtime, crashing the QR
            // scanner only in minified builds. See proguard-rules.pro and
            // https://github.com/juliansteenbakker/mobile_scanner/pull/1726
            isMinifyEnabled = true
            proguardFiles(
                getDefaultProguardFile("proguard-android-optimize.txt"),
                "proguard-rules.pro",
            )
        }
    }

    packaging {
        jniLibs {
            excludes += listOf("**/*.a")
            useLegacyPackaging = true
        }
    }
}

kotlin {
    compilerOptions {
        jvmTarget = org.jetbrains.kotlin.gradle.dsl.JvmTarget.JVM_17
    }
}

flutter {
    source = "../.."
}
