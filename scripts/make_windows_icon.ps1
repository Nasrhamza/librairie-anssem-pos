param(
    [Parameter(Mandatory = $true)]
    [string]$InputPng,

    [Parameter(Mandatory = $true)]
    [string]$OutputIco
)

$ErrorActionPreference = 'Stop'
Add-Type -AssemblyName System.Drawing

$sourcePath = [System.IO.Path]::GetFullPath($InputPng)
$outputPath = [System.IO.Path]::GetFullPath($OutputIco)
$outputDirectory = [System.IO.Path]::GetDirectoryName($outputPath)

if (-not [System.IO.File]::Exists($sourcePath)) {
    throw "Logo source not found: $sourcePath"
}
if (-not [System.IO.Directory]::Exists($outputDirectory)) {
    [System.IO.Directory]::CreateDirectory($outputDirectory) | Out-Null
}

$sizes = @(16, 20, 24, 32, 40, 48, 64, 128, 256)
$frames = New-Object System.Collections.Generic.List[object]
$source = [System.Drawing.Image]::FromFile($sourcePath)

try {
    # The Windows icon uses the clear book-and-pencil mark from the full shop logo.
    $cropX = [int]($source.Width * 0.13)
    $cropY = 0
    $cropWidth = [int]($source.Width * 0.74)
    $cropHeight = [int]($source.Height * 0.72)

    foreach ($size in $sizes) {
        $bitmap = New-Object System.Drawing.Bitmap($size, $size, [System.Drawing.Imaging.PixelFormat]::Format32bppArgb)
        $graphics = [System.Drawing.Graphics]::FromImage($bitmap)
        $memory = New-Object System.IO.MemoryStream
        try {
            $graphics.Clear([System.Drawing.Color]::White)
            $graphics.CompositingQuality = [System.Drawing.Drawing2D.CompositingQuality]::HighQuality
            $graphics.InterpolationMode = [System.Drawing.Drawing2D.InterpolationMode]::HighQualityBicubic
            $graphics.SmoothingMode = [System.Drawing.Drawing2D.SmoothingMode]::HighQuality
            $graphics.PixelOffsetMode = [System.Drawing.Drawing2D.PixelOffsetMode]::HighQuality

            $padding = [Math]::Max(1, [int]($size * 0.07))
            $destination = [System.Drawing.Rectangle]::new(
                $padding,
                $padding,
                ($size - (2 * $padding)),
                ($size - (2 * $padding))
            )
            $graphics.DrawImage(
                $source,
                $destination,
                $cropX,
                $cropY,
                $cropWidth,
                $cropHeight,
                [System.Drawing.GraphicsUnit]::Pixel
            )
            $bitmap.Save($memory, [System.Drawing.Imaging.ImageFormat]::Png)
            $frames.Add([pscustomobject]@{ Size = $size; Data = $memory.ToArray() })
        }
        finally {
            $memory.Dispose()
            $graphics.Dispose()
            $bitmap.Dispose()
        }
    }
}
finally {
    $source.Dispose()
}

$stream = [System.IO.File]::Open($outputPath, [System.IO.FileMode]::Create, [System.IO.FileAccess]::Write)
$writer = New-Object System.IO.BinaryWriter($stream)
try {
    $writer.Write([UInt16]0)
    $writer.Write([UInt16]1)
    $writer.Write([UInt16]$frames.Count)

    $offset = 6 + (16 * $frames.Count)
    foreach ($frame in $frames) {
        $dimension = if ($frame.Size -eq 256) { 0 } else { $frame.Size }
        $writer.Write([Byte]$dimension)
        $writer.Write([Byte]$dimension)
        $writer.Write([Byte]0)
        $writer.Write([Byte]0)
        $writer.Write([UInt16]1)
        $writer.Write([UInt16]32)
        $writer.Write([UInt32]$frame.Data.Length)
        $writer.Write([UInt32]$offset)
        $offset += $frame.Data.Length
    }

    foreach ($frame in $frames) {
        $writer.Write([Byte[]]$frame.Data)
    }
}
finally {
    $writer.Dispose()
    $stream.Dispose()
}

Write-Output "Windows icon created: $outputPath"
