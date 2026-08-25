import 'dart:io';
import 'dart:ui' as ui;
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:qr_flutter/qr_flutter.dart';
import 'package:downloadsfolder/downloadsfolder.dart';
import 'package:path_provider/path_provider.dart';
import 'package:retrovibed/designkit.dart' as ds;

const String _founderSignalGroupInviteUrl =
    "https://signal.group/#CjQKIF533TM1PdI7hJT8ileH9akw1FMMymQnMhn7UjSub1E1EhCJopDCptx_j_9Tu_OJ5_F7";

Future<Uint8List> _renderQrPng(String data, {double size = 1024}) async {
  final painter = QrPainter(
    data: data,
    version: QrVersions.auto,
    dataModuleStyle: QrDataModuleStyle(color: Colors.black, dataModuleShape: QrDataModuleShape.square),
  );
  final recorder = ui.PictureRecorder();
  final canvas = Canvas(recorder);
  canvas.drawRect(Rect.fromLTWH(0, 0, size, size), Paint()..color = Colors.white);
  painter.paint(canvas, Size(size, size));
  final image = await recorder.endRecording().toImage(size.toInt(), size.toInt());
  final bytes = await image.toByteData(format: ui.ImageByteFormat.png);
  return bytes!.buffer.asUint8List(bytes.offsetInBytes, bytes.lengthInBytes);
}

Future<String> _saveToDownloads(Uint8List bytes, String filename) {
  if (Platform.isIOS || Platform.isAndroid) {
    return getTemporaryDirectory().then((dir) {
      final dst = File('${dir.path}/$filename');
      return dst.writeAsBytes(bytes).then((_) {
        return copyFileIntoDownloadFolder(
          dst.path,
          filename,
        ).then(
          (ok) => ok ?? false ? Future.value('Downloads/$filename') : Future.error('failed to copy file $filename'),
        );
      });
    });
  }
  return getDownloadsDirectory().then((dir) async {
    final dst = '${dir!.path}/$filename';
    await File(dst).writeAsBytes(bytes);
    return dst;
  });
}

class SignalGroupDetail extends StatelessWidget {
  final EdgeInsets? margin;
  final EdgeInsets? padding;
  const SignalGroupDetail({super.key, this.margin, this.padding});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final defaults = ds.Defaults.of(context);

    return ds.Container(
      margin: margin ?? defaults.margin,
      padding: padding ?? defaults.padding,
      Column(
        mainAxisSize: MainAxisSize.min,
        spacing: defaults.spacing,
        children: [
          Text("Founder Signal Group", style: theme.textTheme.titleMedium),

          ds.Help(
            ds.Container(
              clipBehavior: Clip.antiAlias,
              constraints: BoxConstraints(
                maxHeight: defaults.compact + defaults.padding.vertical,
              ),
              ClipRRect(
                borderRadius: defaults.borderRadius,
                child: QrImageView(
                  data: _founderSignalGroupInviteUrl,
                  version: QrVersions.auto,
                  backgroundColor: Colors.white,
                  dataModuleStyle: QrDataModuleStyle(color: Colors.black, dataModuleShape: QrDataModuleShape.square),
                ),
              ),
            ),
            ds.Hint(const Text("Scanning this QR code opens the Signal group invite")),
          ),
          Row(
            mainAxisSize: MainAxisSize.min,
            spacing: defaults.spacing,
            children: [
              ds.LoadingIconButton(
                icon: const Icon(Icons.copy),
                tooltip: 'Copy Invite Link',
                onPressed: () => Clipboard.setData(
                  const ClipboardData(text: _founderSignalGroupInviteUrl),
                ),
              ),
              ds.LoadingIconButton(
                icon: const Icon(Icons.download),
                tooltip: 'Download QR Code',
                onPressed: () => _renderQrPng(_founderSignalGroupInviteUrl)
                    .then((bytes) => _saveToDownloads(bytes, 'retro.founder.chat.png'))
                    .then((path) {
                      if (!context.mounted) return;
                      ScaffoldMessenger.of(context).showSnackBar(
                        SnackBar(content: Text('Saved to $path')),
                      );
                    }),
              ),
            ],
          ),
          Text(
            "Scan this QR code to join the founder chat.",
            style: theme.textTheme.bodySmall,
            textAlign: TextAlign.center,
          ),
        ],
      ),
    );
  }
}
