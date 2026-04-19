Pod::Spec.new do |s|
  s.name         = 'RetrovivedBind'
  s.version      = '1.0.0'
  s.summary      = 'Go native bindings'
  s.homepage     = 'https://retrovibe.space'
  s.license      = { :type => 'Proprietary' }
  s.author       = 'retrovibed'
  s.source       = { :path => '.' }
  s.platform     = :ios, '16.0'

  s.vendored_frameworks = 'RetrovivedBind.xcframework'
end
