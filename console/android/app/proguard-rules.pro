# mobile_scanner's bundled consumer rules (proguard-rules.pro, as of 7.2.0)
# use a single-level wildcard (com.google.mlkit.*) which does not match
# nested packages such as com.google.mlkit.vision.barcode.*, letting R8
# strip classes ML Kit resolves at runtime via its component discovery
# service. This causes a NullPointerException in the scanner method
# channel only in minified/release builds. Fixed upstream, not yet
# released: https://github.com/juliansteenbakker/mobile_scanner/pull/1726
# Re-keep the affected packages recursively here until that lands.
-keep class com.google.mlkit.** { *; }
-keep class com.google.android.gms.internal.mlkit** { *; }
-keep class com.google.android.libraries.barhopper.** { *; }
-keep class com.google.photos.** { *; }

-keepclassmembers class * extends java.lang.Enum {
    <fields>;
    public static **[] values();
    public static ** valueOf(java.lang.String);
}
