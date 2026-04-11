Pod::Spec.new do |s|
  s.name         = 'RetrovivedBind'
  s.version      = '1.0.0'
  s.summary      = 'Go native bindings'
  s.homepage     = 'https://retrovibe.space'
  s.license      = { :type => 'Proprietary' }
  s.author       = 'retrovibed'
  s.source       = { :path => '.' }
  s.platform     = :ios, '16.0'

  static_libs = Dir[File.join(__dir__, '*.a')].map { |f| File.basename(f) }
  s.vendored_libraries = static_libs
  s.static_framework = true
  s.source_files = 'Classes/**/*.{h,m}'
  s.public_header_files = 'Classes/**/*.h'
  s.libraries = 'c++', 'resolv'
  s.user_target_xcconfig = {
    'OTHER_LDFLAGS' => '$(inherited) ' + static_libs.map { |f| "-force_load \"${PODS_ROOT}/../#{f}\"" }.join(' ')
  }
end
