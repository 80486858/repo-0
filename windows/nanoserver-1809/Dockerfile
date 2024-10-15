# escape=`

FROM mcr.microsoft.com/powershell:nanoserver-1809 as downloader

ENV VERSION 1.0.6
ENV SHA256 33b95f2cf7387deac27d9f78b768b376919770fc5cc61b9437524f1a1a608c8a

SHELL ["pwsh", "-Command", "$ErrorActionPreference = 'Stop'; $ProgressPreferences = 'SilentlyContinue';"]
RUN New-Item -ItemType Directory -Path C:\download; `
    Invoke-WebRequest `
      -Uri "https://releases.hashicorp.com/consul/${env:VERSION}/consul_${env:VERSION}_windows_amd64.zip" `
      -Outfile C:\download\consul.zip; `
    if ($(Get-FileHash C:\download\consul.zip).Hash -ne ${env:SHA256}) `
    { `
      Write-Error "Error performing checksum."; `
      Exit 1; `
    } `
    Expand-Archive -Path C:\download\consul.zip -DestinationPath C:\download

FROM mcr.microsoft.com/windows/nanoserver:2009-KB4579311
RUN mkdir C:\Consul\Data; `
    mkdir C:\Consul\Config
COPY --from=downloader C:/download/consul.exe C:/Consul/consul.exe

# Expose the consul data directory as a volume since there's mutable state in there.
VOLUME C:\Consul\Data

# Server RPC is used for communication between Consul clients and servers for internal
# request forwarding.
EXPOSE 8300

# Serf LAN and WAN (WAN is used only by Consul servers) are used for gossip between
# Consul agents. LAN is within the datacenter and WAN is between just the Consul
# servers in all datacenters.
EXPOSE 8301 8301/udp 8302 8302/udp

# HTTP and DNS (both TCP and UDP) are the primary interfaces that applications
# use to interact with Consul.
EXPOSE 8500 8600 8600/udp

USER ContainerUser
ENTRYPOINT [ "C:\\Consul\\consul.exe" ]
CMD ["agent", "-ui", "-dev", "-client", "0.0.0.0"]
