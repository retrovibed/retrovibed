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
  s.source_files = 'Classes/**/*.{h,m}'
  s.public_header_files = 'Classes/**/*.h'
  s.libraries = 'c++', 'resolv'
  force_load_flags = static_libs.map { |lib| "-force_load $(PODS_ROOT)/RetrovivedBind/#{lib}" }.join(' ')
  s.pod_target_xcconfig = { 'OTHER_LDFLAGS' => force_load_flags }
end
