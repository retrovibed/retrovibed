import 'dart:io';

import 'package:flutter/material.dart';
import 'package:image_picker/image_picker.dart';
import 'package:permission_handler/permission_handler.dart';
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/design.kit/file.drop.well.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/mimex.dart' as mimex;
import 'api.dart';

typedef FnCapturePhoto = Future<XFile?> Function();
typedef FnCameraPermission = Future<bool> Function();

Future<XFile?> _camera() => ImagePicker().pickImage(source: ImageSource.camera);

// the manifest declares android.permission.CAMERA, so the capture intent needs the
// grant before it will hand anything back.
Future<bool> _cameraPermission() => Permission.camera.request().then((status) => status.isGranted);

// The photo action on a community card: uploads an image into the library and
// immediately publishes it to the community using that community's default
// publish mode. on desktop the image is picked from disk, or dropped onto the
// button. on mobile it is taken with the camera first, then published.
class SocialActionPhoto extends StatelessWidget {
  static const _icon = Icons.add_a_photo_outlined;
  static const _tooltip = "Publish a photo";

  final Community community;
  final media.FnUploadRequest apiupload;
  final FnPublishingPublish apipublish;
  final FnCapturePhoto capture;
  final FnCameraPermission permission;

  const SocialActionPhoto(
    this.community, {
    super.key,
    this.apiupload = media.media.upload,
    this.apipublish = publishing.publish,
    this.capture = _camera,
    this.permission = _cameraPermission,
  });

  // the camera hands back a bare path, so sniff the mimetype the same way
  // FileDropWell does for picked and dropped files.
  Future<String> _mimetype(XFile photo) {
    return File(photo.path)
        .openSync()
        .read(mimex.defaultMagicNumbersMaxLength)
        .then((v) => v.toList())
        .then((bits) => mimex.fromFile(photo.name, magicbits: bits));
  }

  Future<void> _publish(
    List<httpx.Option> auth,
    String path,
    String name,
    String mimetype, {
    ValueNotifier<int>? progress,
  }) {
    return media.media
        .uploadable(path, name, mimetype, progress: progress)
        .then((f) {
          return apiupload((req) {
            req..files.add(f);
            return req;
          });
        })
        .then((resp) {
          final req = PublishContentRequest(
            publishMode: community.defaultPublishMode,
            publishedContent: PublishedContent()
              ..communityId = community.id
              ..knownMediaId = resp.media.knownMediaId
              ..libraryId = resp.media.id,
          );

          return httpx.withRetry(() => apipublish(community.id, req, options: auth));
        });
  }

  @override
  Widget build(BuildContext context) {
    if (ds.Defaults.of(context).mobile) {
      return ds.LoadingIconButton(
        icon: const Icon(_icon),
        tooltip: _tooltip,
        onPressed: () {
          final messenger = ScaffoldMessenger.of(context);

          return permission().then<void>((granted) {
            if (!granted) {
              messenger.showSnackBar(
                const SnackBar(content: Text("camera access is required to publish a photo")),
              );
              return null;
            }

            return capture().then<void>((photo) {
              // the user backed out of the camera without taking anything.
              if (photo == null) return null;
              return _mimetype(photo).then(
                (mimetype) =>
                    _publish([authn.request(authn.AuthzCache.meta(context))], photo.path, photo.name, mimetype),
              );
            });
          });
        },
      );
    }

    return FileDropWell.icon(
      (evt, {progress}) {
        return Future.wait(
          evt.files.map(
            (f) => _publish(
              [authn.request(authn.AuthzCache.meta(context))],
              f.path,
              f.name,
              f.mimeType!,
              progress: progress,
            ),
          ),
        ).then((_) => null);
      },
      icon: _icon,
      tooltip: _tooltip,
      mimetypes: mimex.images,
    );
  }
}
